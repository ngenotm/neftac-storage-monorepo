package policy

import "github.com/casbin/casbin/v2"

var E *casbin.Enforcer

func Init() {
	m, _ := casbin.NewModelFromString(`
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act, eft

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = r.sub == p.sub && regexMatch(r.obj, p.obj) && r.act == p.act
	`)
	E, _ = casbin.NewEnforcer(m)
	E.AddPolicy("admin", ".*", ".*", "allow")
}

func Allow(sub, obj, act string) bool {
	ok, _ := E.Enforce(sub, obj, act)
	return ok
}
