package middleware

import (
	"errors"
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/webserver/db"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

var _c ConfigAuth

const currentUserKey = "current_user"
const currentPassword = "current_password"
const loginFailedKey = "___login_failed"

func AddAuth(authApi, authedApi *gin.RouterGroup) error {
	m, err := newAuthMiddleWare()
	if err != nil {
		return err
	}

	authApi.POST("login", m.LoginHandler)
	authedApi.Use(m.MiddlewareFunc())

	authedApi.POST("auth/logout", m.LogoutHandler)
	return nil
}

func newAuthMiddleWare() (*jwt.GinJWTMiddleware, error) {
	conf, err := GetConfig()
	if err != nil {
		return nil, err
	}
	_c = conf

	c := authConfig(conf)
	return jwt.New(c)
}

func authConfig(conf ConfigAuth) *jwt.GinJWTMiddleware {
	c := jwt.GinJWTMiddleware{
		Realm:           conf.Realm,
		Key:             []byte(conf.Secret),
		Timeout:         time.Duration(conf.Timeout) * time.Second,
		MaxRefresh:      time.Duration(conf.MaxRefresh) * time.Second,
		IdentityKey:     conf.IDKey,
		PayloadFunc:     payloadFunc,
		IdentityHandler: idHandler,
		Authenticator:   authenticator,
		Authorizator:    authorizator,
		Unauthorized:    unAuth,
		LoginResponse:   loginResponse,
		LogoutResponse:  logoutResponse,
		TokenLookup:     "header: Authorization, query: token",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	}
	return &c
}

func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			_c.IDKey: v.JwtKey,
		}
	}
	return jwt.MapClaims{}
}

func idHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	jwtKey := claims[_c.IDKey].(string)

	var user = &User{}
	err := db.DB.Where("jwt_key = ?", jwtKey).First(user).Error
	if err != nil {
		return nil
	}
	return user
}

func CurrentUser(c *gin.Context) *User {
	conf := authConfig(_c)
	user, _ := c.Get(conf.IdentityKey)
	currentUser := user.(*User)
	return currentUser
}

type loginArgs struct {
	Username string `json:"username" binding:"required,alphanum,lte=32"`
	Password string `json:"password" binding:"required,printascii,lte=128"`
}

func authenticator(c *gin.Context) (interface{}, error) {
	var loginValues loginArgs
	if err := c.ShouldBind(&loginValues); err != nil {
		return nil, errors.New("请正确输入用户名密码")
	}

	userName := loginValues.Username
	password := loginValues.Password

	// 直接记录下来， 不管成功与否， 后面看情况使用
	c.Set(loginFailedKey, &User{Username: userName})

	// 密码登录
	var user User
	err := db.DB.Where("username = ?", userName).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) || !user.CheckPassword(password) {
		return nil, errors.New("用户名或者密码填写不对")
	}

	if JwtAuthenticatorHook != nil {
		if err = JwtAuthenticatorHook(c, userName); err != nil {
			return nil, err
		}
	}

	// 注意这个不是idHandler的重复， 这里是给 loginResponse用的
	// 验证的地方有两套流程
	// 1. token 验证， 用到idHandler
	// 2. 登录过程， 此时不会经过中间件，只是单纯的验证密码，发token,成功了调用 loginResponse,失败顿斯 unauthorized
	c.Set(currentUserKey, &user)
	c.Set(currentPassword, password)
	return user, nil

}

func authorizator(data interface{}, c *gin.Context) bool {
	if data == nil {
		return false
	}
	fmt.Println("data", data)
	return true
}

// 密码登录失败时候调用的函数
func unAuth(c *gin.Context, code int, message string) {
	user := &User{}
	if failedUser, ok := c.Get(loginFailedKey); ok {
		user = failedUser.(*User)
	}

	desc := "登录失败"
	msg := message
	if message == "Token is expired" {
		desc = ""
		msg = "您的登录已过期，请重新登录"
	}

	if message == "query token is empty" {
		desc = ""
		msg = "未携带身份凭证"
	}

	IMsg := fmt.Sprintf("用户名：%v，登录失败：%v", user.Username, msg)
	if AuthLog != nil {
		AuthLog("I_LOGIN", "登录", RemoteAddr(c), IMsg, c)
	}
	c.JSON(http.StatusUnauthorized, gin.H{
		"msg":  msg,
		"desc": desc,
	})
}

// 密码登录成功时调用的函数
func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	u, _ := c.Get(currentUserKey)
	user := u.(*User)

	info := user.GetInfo()

	info["token"] = token

	now := time.Now()
	user.LastLoginTime = &now
	db.DB.Save(user)

	if JwtLoginResponseHook != nil {
		password, _ := c.Get(currentPassword)
		JwtLoginResponseHook(user.Username, password.(string), &info)
	}

	msg := fmt.Sprintf("用户名：%v，登录成功", user.Username)
	if AuthLog != nil {
		AuthLog("I_LOGIN", "登录", RemoteAddr(c), msg, c)
	}
	c.JSON(http.StatusOK, info)
}

// 退出登录
func logoutResponse(c *gin.Context, code int) {
	user := CurrentUser(c)
	msg := fmt.Sprintf("用户名：%v，成功退出登录", user.Username)
	if AuthLog != nil {
		AuthLog("I_LOGOUT", "退出登录", RemoteAddr(c), msg, c)
	}
	c.Status(201)
}

func RemoteAddr(c *gin.Context) string {
	addr := c.Request.RemoteAddr
	idx := strings.LastIndex(addr, ":")
	return addr[:idx]
}
