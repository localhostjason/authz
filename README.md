# Authz

1. Casbin 库 权限管理
2. gin router 封装. 用途可初始化路由信息

例子见 `main.go`

````
func SetView(r *gin.Engine) (err error) {

	api := r.Group("api")

    // 加载 casbin 必须已登录成功
	rootGroup := model.NewAuthRootGroup("admin")
	// rootGroup.LoadCasbinConfig(....) 可自定义加载 casbin 配置
	err = rootGroup.LoadCasbin()
	if err != nil {
		return
	}

	g := rootGroup.CreateRootGroup(api)

	AddViewUser(g.Group("用户管理", "user"))

	return
}
````