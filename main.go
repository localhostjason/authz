package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/authz/middleware"
	"github.com/localhostjason/authz/model"
	"github.com/localhostjason/webserver/daemonx"
	"github.com/localhostjason/webserver/db"
	"github.com/localhostjason/webserver/server/util/ue"
	"github.com/localhostjason/webserver/server/util/uv"
	uuid "github.com/satori/go.uuid"
	"time"
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

type User struct {
	middleware.User
	Descx string `json:"descx"`
	Role  string `json:"role" gorm:"type:string;size:64;not null"`
	Desc  string `json:"desc" gorm:"type:string;size:256"`
}

func InitUser() error {
	user := User{
		Role: "admin",
	}
	user.Username = "admin"
	user.JwtKey = uuid.NewV4()
	user.Time = time.Now()
	user.SetPassword("123")
	db.DB.FirstOrCreate(&user)
	return nil
}

func init() {
	db.RegTables(&User{})
	db.AddInitHook(InitUser)
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
	apiAuth := r.Group("api/auth")
	api := r.Group("api")

	jwt := middleware.NewJwt()
	jwt.LoadAuthLog(func(code, action, rip, msg string, c *gin.Context) {
		fmt.Println("login log save to db", code, action, rip, msg)
	})
	//jwt.AuthenticatorHandler(func(c *gin.Context, username string) error {
	//	return errors.New("test error")
	//})
	jwt.LoginResponseHandler(func(username, password string, info *map[string]interface{}) {
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

func main() {
	// 自定义的配置路径 可配置
	s := daemonx.NewMainServer()
	s.LoadView(SetView)
	s.Run()
}
