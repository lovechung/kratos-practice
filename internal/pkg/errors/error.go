package ex

import "github.com/go-kratos/kratos/v2/errors"

var (
	ErrUserNotFound = errors.InternalServer("10000", "该用户不存在")

	ErrCarNotFound = errors.InternalServer("20000", "该汽车不存在")
)
