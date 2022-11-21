# Authz

1. gin router 封装库. 用途可初始化路由信息
2. jwt 封装 （可不引入）

例子见 `main.go`

````
func SetView(r *gin.Engine) (err error) {
	apiAuth := r.Group("api/auth")
	api := r.Group("api")

	jwt := middleware.NewJwt()
	jwt.LoadAuthLog(func(code, action, rip, msg string, c *gin.Context) {
		fmt.Println("login log save to db", code, action, rip, msg)
	})
	//jwt.AuthenticatorHook(func(c *gin.Context, username string) error {
	//	return errors.New("test error")
	//})
	jwt.LoginResponseHook(func(username, password string, info *map[string]interface{}) {
		(*info)["tt"] = "tt"
	})
	err = jwt.AddAuth(apiAuth, api)
	if err != nil {
		return err
	}

	rootGroup := model.NewAuthRootGroup()
	rootGroup.LoadOpLog(func(code, action, rip, msg string, c *gin.Context) {
		fmt.Println("save op log to db:", code, action, rip, msg)
	})

	// 加载 Permissions 必须已登录成功
	rootGroup.LoadPermissionsHandler(func(c *gin.Context) bool {
		return false
	})

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