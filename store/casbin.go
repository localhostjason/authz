package store

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/localhostjason/webserver/db"
)

var E *casbin.Enforcer // 作为全局访问

type CasBin struct {
	ModelType  string
	ConfigFile string
	ConfigText string
}

func NewCasBin(modelType, configFile, configText string) *CasBin {
	if modelType == "" {
		modelType = "text"
	}

	if configText == "" {
		configText = casbinText
	}

	if configFile == "" {
		configFile = GetRBACMConfig()
	}

	return &CasBin{ModelType: modelType, ConfigFile: configFile, ConfigText: configText}
}

func (c *CasBin) Run() error {
	dbx := db.DB
	adapter, err := gormadapter.NewAdapterByDBWithCustomTable(dbx, &CasbinRule{})

	var enforcer *casbin.Enforcer
	if c.ModelType == "text" {
		m, _ := model.NewModelFromString(c.ConfigText)
		enforcer, err = casbin.NewEnforcer(m, adapter)
	} else {
		enforcer, err = casbin.NewEnforcer(c.ConfigFile, adapter)
	}
	err = enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	E = enforcer
	return nil
}

func (c *CasBin) Clear(v int, p ...string) bool {
	ok, _ := E.RemoveFilteredPolicy(v, p...)
	return ok
}

func (c *CasBin) Get() [][]string {
	policy := E.GetPolicy()
	return policy
}

func GetAllPolicy() [][]string {
	policy := E.GetPolicy()
	return policy
}
