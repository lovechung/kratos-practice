package ex

import "github.com/go-kratos/kratos/v2/errors"
import user "github.com/lovechung/api-base/api/user"

var (
	UserNotFound = user.ErrorUserNotFound("该用户不存在")

	CarNotFound = errors.InternalServer("20000", "该汽车不存在")
)
