package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/authz/model"
	"github.com/localhostjason/webserver/daemonx"
	"github.com/localhostjason/webserver/server/util/ue"
	"github.com/localhostjason/webserver/server/util/uv"
)

const (
	I_OP = "I_OP"
)

var iMap = map[string]ue.Info{
	I_OP: {Code: I_OP, Msg: "测试ID: %v, 测试名称：%v"},
}

func init() {
	ue.RegInfos(iMap)
}

func AddViewUserItem(g *model.AuthGroup) {
	g.AddUrl("获取个人信息", model.GET, "info", getU)
}

func AddViewUser(g *model.AuthGroup) {
	AddViewUserItem(g.Group("用户管理", ""))
}

func getU(c *gin.Context) {
	c.Set("OpLog", uv.OP(I_OP, "1", "hello"))
	c.JSON(200, "ok")
}

func SetView(r *gin.Engine) (err error) {
	api := r.Group("api")

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

func main() {
	// 自定义的配置路径 可配置
	s := daemonx.NewMainServer()
	s.LoadView(SetView)
	s.Run()
}
