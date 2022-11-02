package model

import (
	"github.com/gin-gonic/gin"
	"role/store"
)

type AuthRootGroup struct {
	UserRole        string
	CasBinModelType string
	CasBinModelText string
	CasBinModelFile string
}

func NewAuthRootGroup(userRole string) *AuthRootGroup {
	return &AuthRootGroup{UserRole: userRole}
}

// LoadCasbinConfig 不使用默认。去加载配置
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

func (ag *AuthRootGroup) CreateRootGroup(r *gin.RouterGroup) AuthGroup {
	return CreateRootGroup(r, ag.UserRole)
}
