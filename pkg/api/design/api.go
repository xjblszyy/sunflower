package design

import (
	. "goa.design/goa/v3/dsl"
	_ "goa.design/plugins/v3/zaplogger"
)

var _ = API("sunflower", func() {
	Title("微服务")
	HTTP(func() {
		Path("/api")
	})

	Server("sunflower", func() {
		Description("微服务")
		Services("Score")

		Host("localhost", func() {
			Description("default host")
			URI("http://localhost:8000/sunflower")
			URI("grpc://localhost:8080/sunflower")
		})
	})
})

// JWTAuth defines a security scheme that uses JWT tokens.
var JWTAuth = JWTSecurity("jwt", func() {
	Description("使用 JWT 认证, 需要认证的接口添加 ```Header```: ```Authorization: Bearer {jwtToken}```")
	Scope("api:read", "只读权限")
	Scope("api:write", "读写权限")
	Scope("api:admin", "管理权限")
})
