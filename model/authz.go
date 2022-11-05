package model

import (
	"github.com/gin-gonic/gin"
	"github.com/localhostjason/authz/store"
)

var OpLogHook func(code, action, rip, msg string, c *gin.Context)

type AuthRootGroup struct {
	UserRole        string
	CasBinModelType string
	CasBinModelText string
	CasBinModelFile string
}

func NewAuthRootGroup(userRole string) *AuthRootGroup {
	return &AuthRootGroup{UserRole: userRole}
}

// LoadCasbinConfig 不使用默认。可自定义， 默认内置 text config
func (ag *AuthRootGroup) LoadCasbinConfig(modeType, modelText, modelFile string) {
	ag.CasBinModelType = modeType
	ag.CasBinModelText = modelText
	ag.CasBinModelFile = modelFile
}

func (ag *AuthRootGroup) LoadCasbin() error {
	var casBin = store.NewCasBin(ag.CasBinModelType, ag.CasBinModelFile, ag.CasBinModelText)
	if err := casBin.Run(); err != nil {
		return err
	}
	return nil
}

func (ag *AuthRootGroup) LoadOpLog(oplog func(code, action, rip, msg string, c *gin.Context)) {
	OpLogHook = oplog
}

func (ag *AuthRootGroup) CreateRootGroup(r *gin.RouterGroup) AuthGroup {
	return CreateRootGroup(r, ag.UserRole)
}
