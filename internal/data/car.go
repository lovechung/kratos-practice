package data

import (
	"context"
	"entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-practice/internal/biz"
	"kratos-practice/internal/data/ent"
	"kratos-practice/internal/data/ent/car"
	"kratos-practice/internal/data/ent/user"
	ex "kratos-practice/internal/pkg/errors"
	"kratos-practice/internal/pkg/util/pagination"
)

type carRepo struct {
	data *Data
	log  *log.Helper
}

func NewCarRepo(data *Data, logger log.Logger) biz.CarRepo {
	return &carRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r carRepo) ListCar(ctx context.Context,
	page, pageSize int, username *string, model *string) ([]*biz.CarReply, int, error) {
	// 该方法演示多表关联分页多条件查询
	var list []*biz.CarReply

	err := r.data.db.Car.Query().
		Modify(func(s *sql.Selector) {
			selector, t1 := BuildSelector(s, username, model)
			// 分页排序
			selector.
				Offset(pagination.GetOffset(page, pageSize)).
				Limit(pageSize).
				OrderBy(sql.Desc(t1.C(car.FieldRegisteredAt)))
		}).
		Scan(ctx, &list)
	if err != nil || list == nil {
		return nil, 0, err
	}

	total := r.data.db.Car.Query().
		Modify(func(s *sql.Selector) {
			BuildSelector(s, username, model)
		}).
		CountX(ctx)

	return list, total, err
}

// BuildSelector 构建sql，并返回需要后续操作(如:orderBy)的table
func BuildSelector(s *sql.Selector, username *string, model *string) (*sql.Selector, *sql.SelectTable) {
	t1 := sql.Table(car.Table).As("c")
	t2 := sql.Table(user.Table).As("u")

	selector := s.Select(
		t1.C(car.FieldID),
		t1.C(car.FieldModel),
		t1.C(car.FieldRegisteredAt),
		t2.C(user.FieldUsername),
	).
		From(t1).
		LeftJoin(t2).
		On(t1.C(car.FieldUserID), t2.C(user.FieldID))

	// 组装查询条件
	if username != nil {
		selector.Where(sql.P().Contains(user.FieldUsername, *username))
	}
	if model != nil {
		selector.Where(sql.P().Contains(car.FieldModel, *model))
	}
	return selector, t1
}

func (r carRepo) GetById(ctx context.Context, id int64) (*biz.CarReply, error) {
	var rsp []*biz.CarReply
	err := r.data.db.Car.Query().
		Where(car.ID(id)).
		Modify(func(s *sql.Selector) {
			t1 := sql.Table(car.Table).As("c")
			t2 := sql.Table(user.Table).As("u")
			s.Select(
				t1.C(car.FieldID),
				t1.C(car.FieldModel),
				t1.C(car.FieldRegisteredAt),
				t2.C(user.FieldUsername),
			).
				From(t1).
				LeftJoin(t2).
				On(
					t1.C(car.FieldUserID),
					t2.C(user.FieldID),
				)
		}).
		Scan(ctx, &rsp)
	if err != nil {
		return nil, err
	}
	if rsp == nil {
		return nil, ex.ErrCarNotFound
	}
	return rsp[0], err
}

func (r carRepo) Save(ctx context.Context, c *biz.Car) (int64, error) {
	rsp, err := r.data.db.Car.
		Create().
		SetCar(c).
		Save(ctx)
	return rsp.ID, err
}

func (r carRepo) Update(ctx context.Context, c *biz.Car) error {
	return r.data.db.Car.
		Update().
		Where(car.ID(c.ID)).
		SetCar(c).
		Exec(ctx)
}

func (r carRepo) Delete(ctx context.Context, id int64) error {
	return r.data.db.Car.
		DeleteOneID(id).
		Exec(ctx)
}

func ConvertToCar(u *ent.Car) *biz.Car {
	return &biz.Car{
		ID:           u.ID,
		UserID:       &u.UserID,
		Model:        &u.Model,
		RegisteredAt: &u.RegisteredAt,
	}
}
