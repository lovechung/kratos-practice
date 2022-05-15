// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// protoc-gen-go-http v2.2.2

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

type UserHTTPServer interface {
	DeleteUser(context.Context, *DeleteUserReq) (*emptypb.Empty, error)
	GetUser(context.Context, *UserReq) (*UserReply, error)
	ListUser(context.Context, *ListUserReq) (*ListUserReply, error)
	SaveUser(context.Context, *SaveUserReq) (*emptypb.Empty, error)
	UpdateUser(context.Context, *UpdateUserReq) (*emptypb.Empty, error)
}

func RegisterUserHTTPServer(s *http.Server, srv UserHTTPServer) {
	r := s.Route("/")
	r.GET("/user/list", _User_ListUser0_HTTP_Handler(srv))
	r.GET("/user/{id}", _User_GetUser0_HTTP_Handler(srv))
	r.POST("/user", _User_SaveUser0_HTTP_Handler(srv))
	r.PUT("/user", _User_UpdateUser0_HTTP_Handler(srv))
	r.DELETE("/user/{id}", _User_DeleteUser0_HTTP_Handler(srv))
}

func _User_ListUser0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListUserReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.User/ListUser")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListUser(ctx, req.(*ListUserReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListUserReply)
		return ctx.Result(200, reply)
	}
}

func _User_GetUser0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UserReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.User/GetUser")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetUser(ctx, req.(*UserReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*UserReply)
		return ctx.Result(200, reply)
	}
}

func _User_SaveUser0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SaveUserReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.User/SaveUser")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SaveUser(ctx, req.(*SaveUserReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _User_UpdateUser0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateUserReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.User/UpdateUser")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateUser(ctx, req.(*UpdateUserReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _User_DeleteUser0_HTTP_Handler(srv UserHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteUserReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.User/DeleteUser")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteUser(ctx, req.(*DeleteUserReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type UserHTTPClient interface {
	DeleteUser(ctx context.Context, req *DeleteUserReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetUser(ctx context.Context, req *UserReq, opts ...http.CallOption) (rsp *UserReply, err error)
	ListUser(ctx context.Context, req *ListUserReq, opts ...http.CallOption) (rsp *ListUserReply, err error)
	SaveUser(ctx context.Context, req *SaveUserReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateUser(ctx context.Context, req *UpdateUserReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type UserHTTPClientImpl struct {
	cc *http.Client
}

func NewUserHTTPClient(client *http.Client) UserHTTPClient {
	return &UserHTTPClientImpl{client}
}

func (c *UserHTTPClientImpl) DeleteUser(ctx context.Context, in *DeleteUserReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.v1.User/DeleteUser"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) GetUser(ctx context.Context, in *UserReq, opts ...http.CallOption) (*UserReply, error) {
	var out UserReply
	pattern := "/user/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.v1.User/GetUser"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) ListUser(ctx context.Context, in *ListUserReq, opts ...http.CallOption) (*ListUserReply, error) {
	var out ListUserReply
	pattern := "/user/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.v1.User/ListUser"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) SaveUser(ctx context.Context, in *SaveUserReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.v1.User/SaveUser"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *UserHTTPClientImpl) UpdateUser(ctx context.Context, in *UpdateUserReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/user"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.v1.User/UpdateUser"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

type CarHTTPServer interface {
	DeleteCar(context.Context, *DeleteCarReq) (*emptypb.Empty, error)
	GetCar(context.Context, *CarReq) (*CarReply, error)
	ListCar(context.Context, *ListCarReq) (*ListCarReply, error)
	SaveCar(context.Context, *SaveCarReq) (*emptypb.Empty, error)
	UpdateCar(context.Context, *UpdateCarReq) (*emptypb.Empty, error)
}

func RegisterCarHTTPServer(s *http.Server, srv CarHTTPServer) {
	r := s.Route("/")
	r.GET("/car/list", _Car_ListCar0_HTTP_Handler(srv))
	r.GET("/car/{id}", _Car_GetCar0_HTTP_Handler(srv))
	r.POST("/car", _Car_SaveCar0_HTTP_Handler(srv))
	r.PUT("/car", _Car_UpdateCar0_HTTP_Handler(srv))
	r.DELETE("/car/{id}", _Car_DeleteCar0_HTTP_Handler(srv))
}

func _Car_ListCar0_HTTP_Handler(srv CarHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListCarReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.Car/ListCar")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListCar(ctx, req.(*ListCarReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListCarReply)
		return ctx.Result(200, reply)
	}
}

func _Car_GetCar0_HTTP_Handler(srv CarHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CarReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.Car/GetCar")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetCar(ctx, req.(*CarReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*CarReply)
		return ctx.Result(200, reply)
	}
}

func _Car_SaveCar0_HTTP_Handler(srv CarHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in SaveCarReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.Car/SaveCar")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SaveCar(ctx, req.(*SaveCarReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Car_UpdateCar0_HTTP_Handler(srv CarHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in UpdateCarReq
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.Car/UpdateCar")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateCar(ctx, req.(*UpdateCarReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

func _Car_DeleteCar0_HTTP_Handler(srv CarHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in DeleteCarReq
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		if err := ctx.BindVars(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, "/api.v1.Car/DeleteCar")
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.DeleteCar(ctx, req.(*DeleteCarReq))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*emptypb.Empty)
		return ctx.Result(200, reply)
	}
}

type CarHTTPClient interface {
	DeleteCar(ctx context.Context, req *DeleteCarReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	GetCar(ctx context.Context, req *CarReq, opts ...http.CallOption) (rsp *CarReply, err error)
	ListCar(ctx context.Context, req *ListCarReq, opts ...http.CallOption) (rsp *ListCarReply, err error)
	SaveCar(ctx context.Context, req *SaveCarReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
	UpdateCar(ctx context.Context, req *UpdateCarReq, opts ...http.CallOption) (rsp *emptypb.Empty, err error)
}

type CarHTTPClientImpl struct {
	cc *http.Client
}

func NewCarHTTPClient(client *http.Client) CarHTTPClient {
	return &CarHTTPClientImpl{client}
}

func (c *CarHTTPClientImpl) DeleteCar(ctx context.Context, in *DeleteCarReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/car/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.v1.Car/DeleteCar"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "DELETE", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CarHTTPClientImpl) GetCar(ctx context.Context, in *CarReq, opts ...http.CallOption) (*CarReply, error) {
	var out CarReply
	pattern := "/car/{id}"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.v1.Car/GetCar"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CarHTTPClientImpl) ListCar(ctx context.Context, in *ListCarReq, opts ...http.CallOption) (*ListCarReply, error) {
	var out ListCarReply
	pattern := "/car/list"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation("/api.v1.Car/ListCar"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CarHTTPClientImpl) SaveCar(ctx context.Context, in *SaveCarReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/car"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.v1.Car/SaveCar"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}

func (c *CarHTTPClientImpl) UpdateCar(ctx context.Context, in *UpdateCarReq, opts ...http.CallOption) (*emptypb.Empty, error) {
	var out emptypb.Empty
	pattern := "/car"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation("/api.v1.Car/UpdateCar"))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "PUT", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
