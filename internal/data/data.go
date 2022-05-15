package data

import (
	"context"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/nitishm/go-rejson/v4"
	"kratos-practice/internal/biz"
	"kratos-practice/internal/conf"
	"kratos-practice/internal/data/ent"
)

var ProviderSet = wire.NewSet(NewTransaction, NewData, NewDB, NewRedis, NewUserRepo, NewCarRepo)

type Data struct {
	db     *ent.Client
	rdb    *redis.Client
	rejson *rejson.Handler
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

func NewData(db *ent.Client, rdb *redis.Client, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}

	// redis扩展 re-json
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClient(rdb)

	return &Data{
		db:     db,
		rdb:    rdb,
		rejson: rh,
	}, cleanup, nil
}

func NewDB(conf *conf.Data, logger log.Logger) *ent.Client {
	thisLog := log.NewHelper(logger)

	//db, err := ent.Open(
	//	conf.Database.Driver,
	//	conf.Database.Source,
	//)

	drv, err := sql.Open(
		conf.Database.Driver,
		conf.Database.Source,
	)
	sqlDrv := dialect.DebugWithContext(drv, func(ctx context.Context, i ...interface{}) {
		thisLog.WithContext(ctx).Info(i...)
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

	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), conf.Redis.DialTimeout.AsDuration())
	defer cancelFunc()
	err := rdb.Ping(timeout).Err()
	if err != nil {
		thisLog.Fatalf("redis连接失败: %v", err)
	}
	return rdb
}
