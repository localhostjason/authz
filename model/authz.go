package model

import (
	"github.com/gin-gonic/gin"
)

var OpLogHook func(code, action, rip, msg string, c *gin.Context)

type AuthRootGroup struct {
}

func NewAuthRootGroup() *AuthRootGroup {
	return &AuthRootGroup{}
}

func (ag *AuthRootGroup) LoadOpLog(oplog func(code, action, rip, msg string, c *gin.Context)) {
	OpLogHook = oplog
}

func (ag *AuthRootGroup) CreateRootGroup(r *gin.RouterGroup) AuthGroup {
	return CreateRootGroup(r)
}

var PermissionsHandler func(c *gin.Context) bool

func (ag *AuthRootGroup) LoadPermissionsHandler(permissionsHandler func(c *gin.Context) bool) {
	PermissionsHandler = permissionsHandler
}
