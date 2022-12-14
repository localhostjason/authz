package model

import (
	"github.com/gin-gonic/gin"
	. "github.com/localhostjason/authz/store"
	"github.com/localhostjason/webserver/db"
	"github.com/localhostjason/webserver/server/util/ue"
	log "github.com/sirupsen/logrus"
	"net/http"
	"path"
	"runtime/debug"
	"strings"
	"sync"
)

const (
	GET    = "GET"
	PUT    = "PUT"
	POST   = "POST"
	DELETE = "DELETE"
)

const _rootGroupName = "root"

func CreateRootGroup(r *gin.RouterGroup) AuthGroup {
	g := AuthGroup{
		Name:        _rootGroupName,
		Url:         "",
		RouterGroup: r,
		SubGroup:    nil,
		Parent:      nil,
	}
	return g
}

type AuthGroup struct {
	Name        string
	Url         string
	RouterGroup *gin.RouterGroup
	SubGroup    []AuthGroup
	Parent      *AuthGroup
}

func (ag *AuthGroup) Group(name, url string) *AuthGroup {
	r := ag.RouterGroup
	if url != "" {
		r = r.Group(url)
	}

	newGroup := AuthGroup{
		Name:        name,
		Url:         url,
		RouterGroup: r,
		SubGroup:    nil,
		Parent:      ag,
	}
	ag.SubGroup = append(ag.SubGroup, newGroup)
	return &newGroup
}

func (ag *AuthGroup) getPermissionGroup() string {
	var name = ag.Name
	group := ag.Parent
	for group.Name != _rootGroupName {
		name = group.Name
		group = group.Parent
	}
	return name
}

func (ag *AuthGroup) addHandler(method string, url string, h gin.HandlerFunc) {
	switch method {
	case GET:
		ag.RouterGroup.GET(url, h)
	case POST:
		ag.RouterGroup.POST(url, h)
	case PUT:
		ag.RouterGroup.PUT(url, h)
	case DELETE:
		ag.RouterGroup.DELETE(url, h)
	default:
		panic("unknown method")
	}
}

func (ag *AuthGroup) AddUrl(name, method, url string, h gin.HandlerFunc, permissionRoles ...string) {
	GroupName := ag.getPermissionGroup()
	handler := permissionHandler(h, name, method)
	ag.addHandler(method, url, handler)

	fullUrl := path.Join(ag.RouterGroup.BasePath(), url)
	initAuthz(GroupName, name, fullUrl, method)
}

func permissionHandler(h gin.HandlerFunc, action, method string) gin.HandlerFunc {
	var lock sync.Mutex
	return func(c *gin.Context) {
		defer func() {
			msg := ""
			if r := recover(); r != nil {
				if err, ok := r.(*ue.Error); ok {
					msg = err.Error()
					c.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
				} else {
					msg = "????????????"
					log.Error(string(debug.Stack()))
					c.AbortWithStatus(http.StatusInternalServerError)
				}
			}

			logMsg, _ := c.Get("OpLog")
			if OpLogHook != nil {
				if info, ok := logMsg.(*ue.Info); ok {
					if msg == "" {
						msg = info.Msg
					}

					OpLogHook(info.Code, action, RemoteAddr(c), info.Msg, c)
				}
			}
		}()

		if PermissionsHandler != nil && PermissionsHandler(c) {
			c.AbortWithStatus(403)
		} else {
			lock.Lock()
			defer lock.Unlock()
			h(c)
		}
	}
}

// ???????????????
// api ?????? ???????????????acl, ???????????? ???????????????????????? authz ?????????????????????api??????
// casbin ?????????????????????
func initAuthz(groupName, apiName, url, method string) {
	var authz []Authz
	z := Authz{GroupName: groupName, ApiName: apiName, Url: url, Method: method}
	db.DB.Limit(1).Find(&authz, z)

	if len(authz) == 0 {
		db.DB.Create(&z)
	}
}

func RemoteAddr(c *gin.Context) string {
	addr := c.Request.RemoteAddr
	idx := strings.LastIndex(addr, ":")
	return addr[:idx]
}
