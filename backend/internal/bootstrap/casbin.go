package bootstrap

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin(db *gorm.DB) (*casbin.Enforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	m, err := model.NewModelFromString(`
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == '*')
`)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// 添加默认策略
	_, _ = enforcer.AddPolicy("superadmin", "/api/v1/admin/*", "*")
	_, _ = enforcer.AddPolicy("manager", "/api/v1/admin/*", "GET")
	_, _ = enforcer.AddPolicy("manager", "/api/v1/admin/*", "POST")
	_, _ = enforcer.AddPolicy("manager", "/api/v1/admin/*", "PUT")
	_, _ = enforcer.AddPolicy("staff", "/api/v1/admin/dashboard", "GET")
	_, _ = enforcer.AddPolicy("staff", "/api/v1/admin/orders", "GET")
	_, _ = enforcer.AddPolicy("staff", "/api/v1/admin/orders/*", "PUT")
	_, _ = enforcer.AddPolicy("staff", "/api/v1/admin/reviews", "GET")
	_, _ = enforcer.AddPolicy("staff", "/api/v1/admin/reviews/*", "POST")

	Enforcer = enforcer
	return enforcer, nil
}
