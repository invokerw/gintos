package errs

import "github/invokerw/gintos/common/resp"

// DB Err 300000 - 399999

var (
	ErrUserNotFound      = resp.NewErr(300000, "user not found")
	ErrUserPasswordWrong = resp.NewErr(300001, "password wrong")
	ErrTokenExpired      = resp.NewErr(300002, "token expired")
	ErrAvatarDataWrong   = resp.NewErr(300003, "avatar data wrong")
	ErrAvatarExtWrong    = resp.NewErr(300004, "avatar ext wrong")
)
