package biz

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
)

func NewCasbinEnforcer(adapter persist.Adapter) (*casbin.Enforcer, error) {
	model := `
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
	return casbin.NewEnforcer(model, adapter)
}
