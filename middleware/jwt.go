package middleware

import (
	"github.com/gin-gonic/gin"
)

var AuthLog func(arg ...interface{})

type loginResponseHookFunc func(username, password string, info *map[string]interface{})
type authValidateHookFunc func(c *gin.Context, username string) error

var JwtAuthenticatorHook authValidateHookFunc
var JwtLoginResponseHook loginResponseHookFunc

type Jwt struct {
}

func NewJwt() *Jwt {
	return &Jwt{}
}

func (j *Jwt) AddLog(authLog func(arg ...interface{})) {
	AuthLog = authLog
}

// AuthenticatorHook 登录前 可以 自定义 函数 在登录前做限制，ip限制
func (j *Jwt) AuthenticatorHook(authFunc authValidateHookFunc) {
	JwtAuthenticatorHook = authFunc
}

// LoginResponseHook 登录成功后，可以 自定义 函数，在登录成功后做一些事情：比如 记录历史密码，更新密码强度
func (j *Jwt) LoginResponseHook(loginResponse loginResponseHookFunc) {
	JwtLoginResponseHook = loginResponse
}

func (j *Jwt) AddAuth(authApi, authedApi *gin.RouterGroup) error {
	return AddAuth(authApi, authedApi)
}
