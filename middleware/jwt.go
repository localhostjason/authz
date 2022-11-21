package middleware

import (
	"github.com/gin-gonic/gin"
)

var AuthLog func(code, action, rip, msg string, c *gin.Context)

type loginResponseHandlerFunc func(username, password string, info *map[string]interface{})
type authValidateHandlerFunc func(c *gin.Context, username string) error

var JwtAuthenticatorHandler authValidateHandlerFunc
var JwtLoginResponseHandler loginResponseHandlerFunc

type Jwt struct{}

func NewJwt() *Jwt {
	return &Jwt{}
}

func (j *Jwt) LoadAuthLog(authLog func(code, action, rip, msg string, c *gin.Context)) {
	AuthLog = authLog
}

// AuthenticatorHandler 登录前 可以 自定义 函数 在登录前做限制，ip限制
func (j *Jwt) AuthenticatorHandler(authFunc authValidateHandlerFunc) {
	JwtAuthenticatorHandler = authFunc
}

// LoginResponseHandler 登录成功后，可以 自定义 函数，在登录成功后做一些事情：比如 记录历史密码，更新密码强度
func (j *Jwt) LoginResponseHandler(loginResponse loginResponseHandlerFunc) {
	JwtLoginResponseHandler = loginResponse
}

func (j *Jwt) AddAuth(authApi, authedApi *gin.RouterGroup) error {
	return AddAuth(authApi, authedApi)
}
