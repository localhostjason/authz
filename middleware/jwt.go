package middleware

import "github.com/gin-gonic/gin"

var AuthLog func(arg ...interface{})

type Jwt struct {
}

func NewJwt() *Jwt {
	return &Jwt{}
}

func (j *Jwt) AddLog(authLog func(arg ...interface{})) {
	AuthLog = authLog
}

func (j *Jwt) AddAuth(authApi, authedApi *gin.RouterGroup) error {
	return AddAuth(authApi, authedApi)
}
