package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Car struct {
	ID           int64
	UserID       *int64
	Model        *string
	RegisteredAt *time.Time
}

type CarReply struct {
	Id           int64
	Model        string
	RegisteredAt time.Time `sql:"registered_at"`
	UserName     string
}

type CarRepo interface {
	ListCar(ctx context.Context, page, pageSize int, username *string, model *string) ([]*CarReply, int, error)
	GetById(ctx context.Context, id int64) (*CarReply, error)
	Save(context.Context, *Car) (int64, error)
	Update(context.Context, *Car) error
	Delete(ctx context.Context, id int64) error
}

type CarUseCase struct {
	r   CarRepo
	log *log.Helper
	tx  Transaction
}

func NewCarUseCase(r CarRepo, tx Transaction, logger log.Logger) *CarUseCase {
	return &CarUseCase{r: r, tx: tx, log: log.NewHelper(logger)}
}

func (uc *CarUseCase) ListCar(ctx context.Context,
	page, pageSize int, username *string, model *string) ([]*CarReply, int, error) {
	return uc.r.ListCar(ctx, page, pageSize, username, model)
}

func (uc *CarUseCase) GetCarById(ctx context.Context, id int64) (*CarReply, error) {
	return uc.r.GetById(ctx, id)
}

func (uc *CarUseCase) SaveCar(ctx context.Context, c *Car) error {
	_, err := uc.r.Save(ctx, c)
	return err
}

func (uc *CarUseCase) UpdateCar(ctx context.Context, c *Car) error {
	return uc.r.Update(ctx, c)
}

func (uc *CarUseCase) DeleteCar(ctx context.Context, id int64) error {
	return uc.r.Delete(ctx, id)
}
