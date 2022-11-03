package store

import (
	"errors"
	"github.com/localhostjason/webserver/db"
	"gorm.io/gorm"
)

type CasbinRule struct {
	ID        uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	PType     string `json:"p_type" gorm:"column:ptype" description:"策略类型"`
	Role      string `json:"role" gorm:"column:v0" description:"角色"`
	Path      string `json:"path" gorm:"column:v1" description:"api路径 obj"`
	Method    string `json:"method" gorm:"column:v2" description:"访问方法 act"`
	V3        string `gorm:"column:v3"`
	V4        string `gorm:"column:v4"`
	V5        string `gorm:"column:v5" `
	V6        string `gorm:"column:v6" `
	V7        string `gorm:"column:v7" `
	ApiName   string `json:"-"`
	GroupName string `json:"-"`
	Desc      string `json:"desc" description:"策略描述"`
}

func GetCasBins() []CasbinRule {
	var casbinRule []CasbinRule
	db.DB.Find(&casbinRule)
	for i, _ := range casbinRule {
		apiName, groupName := GetApiName(casbinRule[i].Path)
		casbinRule[i].ApiName = apiName
		casbinRule[i].GroupName = groupName
	}
	return casbinRule
}

func GetApiName(path string) (string, string) {
	var authz Authz
	db.DB.Where("url = ?", path).First(&authz)
	return authz.ApiName, authz.GroupName
}

func (c *CasbinRule) Create() error {
	if success, _ := E.AddPolicy(c.Role, c.Path, c.Method); !success {
		return errors.New("存在相同的策略，添加失败")
	}
	return nil
}

func (c *CasbinRule) UpdateDesc(id uint, desc string) error {
	var rule CasbinRule
	err := db.DB.First(&rule, &CasbinRule{ID: id}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("未找到策略")
	}

	rule.Desc = desc
	return db.DB.Save(&rule).Error
}

func (c *CasbinRule) Update(role, path, method string) error {
	updated, err := E.UpdatePolicy([]string{c.Role, c.Path, c.Method}, []string{role, path, method})
	if err != nil {
		return err
	}
	if !updated {
		return errors.New("更新策略失败")
	}
	return nil
}

func (c *CasbinRule) Delete() error {
	ok, err := E.RemoveFilteredNamedPolicy(c.PType, 0, c.Role, c.Path, c.Method)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("已删除策略")
	}
	return nil
}

func (c *CasbinRule) List() [][]string {
	policy := E.GetFilteredPolicy(0, c.Role)
	return policy
}

type Authz struct {
	ID        int    `json:"id"`
	GroupName string `json:"group_name"`
	ApiName   string `json:"api_name"`
	Url       string `json:"url" description:"对象 obj"`
	Method    string `json:"method" description:"act"`
}

func init() {
	db.RegTables(&Authz{})
}
