package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

var (
	ErrUserNotFound = errors.InternalServer("10000", "该用户不存在")
)

type User struct {
	Id        int64
	Username  *string
	Password  *string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type UserRepo interface {
	ListUser(ctx context.Context, page, pageSize int, username *string) ([]*User, int, error)
	GetById(ctx context.Context, id int64) (*User, error)
	Save(context.Context, *User) (int64, error)
	Update(context.Context, *User) error
	Delete(ctx context.Context, id int64) error
}

type UserUseCase struct {
	r   UserRepo
	log *log.Helper
	tx  Transaction
}

func NewUserUseCase(r UserRepo, tx Transaction, logger log.Logger) *UserUseCase {
	return &UserUseCase{r: r, tx: tx, log: log.NewHelper(logger)}
}

func (uc *UserUseCase) ListUser(ctx context.Context, page, pageSize int, username *string) ([]*User, int, error) {
	return uc.r.ListUser(ctx, page, pageSize, username)
}

func (uc *UserUseCase) GetUserById(ctx context.Context, id int64) (*User, error) {
	return uc.r.GetById(ctx, id)
}

func (uc *UserUseCase) SaveUser(ctx context.Context, u *User) error {
	id, err := uc.r.Save(ctx, u)
	uc.log.Infof("新增的用户id=%d", id)
	return err
}

func (uc *UserUseCase) UpdateUser(ctx context.Context, u *User) error {
	var err error
	// 带有事务的操作
	if e := uc.tx.ExecTx(ctx, func(ctx context.Context) error {
		err = uc.r.Update(ctx, u)
		if err != nil {
			return err
		}
		return nil
	}); e != nil {
		return err
	}
	return nil
}

func (uc *UserUseCase) DeleteUser(ctx context.Context, id int64) error {
	return uc.r.Delete(ctx, id)
}
