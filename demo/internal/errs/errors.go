package errs

import "github/invokerw/gintos/common/resp"

// DB Err 300000 - 399999

var (
	ErrUserNotFound      = resp.NewErr(300000, "user not found")
	ErrUserPasswordWrong = resp.NewErr(300001, "password wrong")
	ErrTokenExpired      = resp.NewErr(300002, "token expired")
)
