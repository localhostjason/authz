package store

import (
	"github.com/localhostjason/webserver/db"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCasbinAdd(t *testing.T) {
	assert.Nil(t, db.Connect())

	casbin := NewCasBin("", "", "")
	assert.Nil(t, casbin.Run())

	c := CasbinRule{
		Path:   "/api/user/info",
		Method: "PUT",
		Role:   "admin",
	}

	err := c.Create()
	assert.Nil(t, err)
	t.Log(GetAllPolicy())
}

func TestCasbinRemove(t *testing.T) {
	assert.Nil(t, db.Connect())

	casbin := NewCasBin("", "", "")
	assert.Nil(t, casbin.Run())

	c := CasbinRule{
		PType:  "p",
		Path:   "/api/user/info",
		Method: "PUT",
		Role:   "admin",
	}

	_ = c.Create()
	t.Log(GetAllPolicy())

	err := c.Delete()
	assert.Nil(t, err)
	t.Log(GetAllPolicy())
}

func TestCasbinUpdate(t *testing.T) {
	assert.Nil(t, db.Connect())

	casbin := NewCasBin("", "", "")
	assert.Nil(t, casbin.Run())

	c := CasbinRule{
		PType:  "p",
		Path:   "/api/user/info",
		Method: "PUT",
		Role:   "admin",
	}

	_ = c.Create()
	t.Log(GetAllPolicy())

	err := c.Update("test", "/api/user/info", "PUT")
	assert.Nil(t, err)
	t.Log(GetAllPolicy())
}

func TestCasbinUpdateApiName(t *testing.T) {
	assert.Nil(t, db.Connect())

	casbin := NewCasBin("", "", "")
	assert.Nil(t, casbin.Run())

	c := CasbinRule{
		PType:  "p",
		Path:   "/api/user/info",
		Method: "PUT",
		Role:   "admin",
	}

	_ = c.Create()
	t.Log(GetAllPolicy())

	err := c.UpdateApiName("测试")
	assert.Nil(t, err)
	t.Log(GetAllPolicy())
}
