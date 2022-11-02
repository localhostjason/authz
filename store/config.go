package store

import "github.com/localhostjason/webserver/server/config"

const _key = "rbac_model"

func GetRBACMConfig() string {
	var c string
	_ = config.GetConfig(_key, &c)
	return c
}

func init() {
	_ = config.RegConfig(_key, "/tmp/rbac_models.conf")
}
