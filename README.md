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

封装了一套简单的jwt,可不使用

```golang

    jwt := middleware.NewJwt()
    
    // 日志 函数 hook
	jwt.AddLog(func(arg ...interface{}) {
		fmt.Println(111, arg)
	})
    // 登录前的 验证 hook
	//jwt.AuthenticatorHook(func(c *gin.Context, username string) error {
	//	return errors.New("test error")
	//})
	
	// 登录成功后 函数
	jwt.LoginResponseHook(func(username, password string, info *map[string]interface{}) {
		(*info)["tt"] = "tt"
	})
```