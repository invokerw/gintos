package biz

import (
	"github/invokerw/gintos/demo/internal/pkg/casbin_logger"
	"github/invokerw/gintos/log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

func NewCasbinEnforcer(adapter persist.Adapter, logger log.Logger) (*casbin.Enforcer, error) {
	mStr := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	m, err := model.NewModelFromString(mStr)
	if err != nil {
		return nil, err
	}
	l := casbin_logger.NewCasbinLogger(log.NewHelper(log.With(logger, "module", "casbin")))
	l.EnableLog(true)
	return casbin.NewEnforcer(m, adapter, l)
}
