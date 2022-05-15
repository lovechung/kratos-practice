package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"kratos-practice/internal/pkg/util/pagination"
	"kratos-practice/internal/pkg/util/time"

	"kratos-practice/api/v1"
	"kratos-practice/internal/biz"
)

type UserService struct {
	v1.UnimplementedUserServer

	uc  *biz.UserUseCase
	log *log.Helper
}

func NewUserService(uc *biz.UserUseCase, logger log.Logger) *UserService {
	return &UserService{uc: uc, log: log.NewHelper(logger)}
}

func (s *UserService) ListUser(ctx context.Context, req *v1.ListUserReq) (*v1.ListUserReply, error) {
	page, pageSize := pagination.GetPage(req.Page, req.PageSize)
	list, total, err := s.uc.ListUser(ctx, page, pageSize, req.Username)

	rsp := &v1.ListUserReply{}
	rsp.Total = int32(total)
	for _, user := range list {
		userInfo := ConvertToUserReply(user)
		rsp.List = append(rsp.List, userInfo)
	}
	return rsp, err
}

func (s *UserService) GetUser(ctx context.Context, req *v1.UserReq) (*v1.UserReply, error) {
	user, err := s.uc.GetUserById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return ConvertToUserReply(user), err
}

func ConvertToUserReply(u *biz.User) *v1.UserReply {
	return &v1.UserReply{
		Id:        u.Id,
		Username:  *u.Username,
		Password:  *u.Password,
		CreatedAt: time.Format(*u.CreatedAt),
	}
}

func (s *UserService) SaveUser(ctx context.Context, req *v1.SaveUserReq) (*emptypb.Empty, error) {
	err := s.uc.SaveUser(ctx, &biz.User{
		Username: &req.Username,
		Password: &req.Password,
	})
	return nil, err
}

func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserReq) (*emptypb.Empty, error) {
	err := s.uc.UpdateUser(ctx, &biz.User{
		Id:       req.Id,
		Username: req.Username,
		Password: req.Password,
	})
	return nil, err
}

func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserReq) (*emptypb.Empty, error) {
	err := s.uc.DeleteUser(ctx, req.Id)
	return nil, err
}
