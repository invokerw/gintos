package errs

import "github/invokerw/gintos/common/resp"

// DB Err 100000 - 199999

var (
	DBErrEntError     = resp.NewErr(100000, "ent error")
	DBErrInvalidParam = resp.NewErr(100001, "invalid param")
	DBErrUserNotFound = resp.NewErr(100010, "user not found")
)
