package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/authz/middleware"
	"github.com/localhostjason/authz/model"
	"github.com/localhostjason/authz/store"
	"github.com/localhostjason/webserver/daemonx"
	"github.com/localhostjason/webserver/db"
	uuid "github.com/satori/go.uuid"
	"time"
)

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

func addRole(c *gin.Context) {
	r := store.CasbinRule{
		Path:   "/api/user/info",
		Method: "GET",
		Role:   "admin",
	}

	err := r.Create()
	if err != nil {
		fmt.Println(11, err)
	}
	c.Status(201)
}

func deleteRole(c *gin.Context) {
	var casbin store.CasbinRule
	err := casbin.Delete(1)
	fmt.Println(err)
	c.JSON(200, store.GetAllPolicy())
}

func updateRole(c *gin.Context) {
	var casbin store.CasbinRule
	err := casbin.Update(1, "admin", "/api/user/password", "GET")
	fmt.Println(err)
	c.JSON(200, store.GetAllPolicy())
}

func AddViewUserItem(g *model.AuthGroup) {
	g.AddUrl("获取个人信息", model.GET, "info", getU)
	g.AddUrl("更改个人密码", model.GET, "password", getU)
	g.AddUrl("增加role", model.GET, "role", addRole)
	g.AddUrl("删除role", model.GET, "del_role", deleteRole)
	g.AddUrl("更新role", model.GET, "update_role", updateRole)
}

func AddViewUser(g *model.AuthGroup) {
	AddViewUserItem(g.Group("用户管理", ""))
}

func getU(c *gin.Context) {
	data := map[string]interface{}{
		"policy":    store.GetAllPolicy(),
		"db_policy": store.GetCasBins(),
	}

	c.JSON(200, data)
}

func SetView(r *gin.Engine) (err error) {
	apiAuth := r.Group("api/auth")
	api := r.Group("api")

	jwt := middleware.NewJwt()
	jwt.AddLog(func(arg ...interface{}) {
		fmt.Println(111, arg)
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

	// 加载 casbin 必须已登录成功
	rootGroup := model.NewAuthRootGroup("admin")
	err = rootGroup.LoadCasbin()
	if err != nil {
		return
	}

	g := rootGroup.CreateRootGroup(api)

	AddViewUser(g.Group("用户管理", "user"))

	return
}

func main() {
	// 自定义的配置路径 可配置
	const defaultConfigPath = "D:\\center\\console\\console.json"
	s := daemonx.NewMainServer(defaultConfigPath, SetView)
	s.Run()
}
