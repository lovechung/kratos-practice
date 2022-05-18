package data

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/nitishm/go-rejson/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"kratos-practice/internal/biz"
	"kratos-practice/internal/conf"
	"kratos-practice/internal/data/ent"
	"time"
)

var ProviderSet = wire.NewSet(NewTransaction, NewData, NewDB, NewRedis, NewUserRepo, NewCarRepo)

type Data struct {
	db  *ent.Client
	rds *redis.Client
	rjs *rejson.Handler
}

type contextTxKey struct{}

func (d *Data) ExecTx(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := d.db.Tx(ctx)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, contextTxKey{}, tx)
	if err := f(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (d *Data) User(ctx context.Context) *ent.UserClient {
	tx, ok := ctx.Value(contextTxKey{}).(*ent.Tx)
	if ok {
		return tx.User
	}
	return d.db.User
}

func NewTransaction(d *Data) biz.Transaction {
	return d
}

func NewData(db *ent.Client, rds *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
		if err := db.Close(); err != nil {
			log.Error(err)
		}
		if err := rds.Close(); err != nil {
			log.Error(err)
		}
	}

	// redis扩展 rejson
	rjs := rejson.NewReJSONHandler()
	rjs.SetGoRedisClient(rds)

	return &Data{
		db:  db,
		rds: rds,
		rjs: rjs,
	}, cleanup, nil
}

func NewDB(conf *conf.Data, logger log.Logger) *ent.Client {
	thisLog := log.NewHelper(logger)

	drv, err := sql.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	// 打印sql日志
	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		thisLog.WithContext(ctx).Debug(i...)
		// 开启db trace
		tracer := otel.Tracer("entgo.io/ent")
		_, span := tracer.Start(ctx,
			"query",
			trace.WithAttributes(
				attribute.String("sql", fmt.Sprint(i...)),
			),
		)
		defer span.End()
	})
	db := ent.NewClient(ent.Driver(sqlDrv))

	if err != nil {
		thisLog.Fatalf("数据库连接失败: %v", err)
	}
	// 运行自动创建表
	//if err := db.Schema.Create(context.Background(), migrate.WithForeignKeys(false)); err != nil {
	//	thisLog.Fatalf("创建表失败: %v", err)
	//}
	return db
}

func NewRedis(conf *conf.Data, logger log.Logger) *redis.Client {
	thisLog := log.NewHelper(logger)

	rds := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	// 开启redis trace
	rds.AddHook(redisotel.TracingHook{})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	err := rds.Ping(timeout).Err()
	if err != nil {
		thisLog.Fatalf("redis连接失败: %v", err)
	}
	return rds
}
