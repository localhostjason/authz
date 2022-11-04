package middleware

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User 要用authz jwt 继承过去
type User struct {
	Id            int64      `json:"id" gorm:"primaryKey"`
	Username      string     `json:"username" gorm:"type:string;size:64;unique;not null"`
	Password      string     `json:"password,omitempty" gorm:"column:_password;type:string;size:128"`
	LastLoginTime *time.Time `json:"last_login_time"`
	Time          time.Time  `json:"time"`                           // 创建时间
	JwtKey        uuid.UUID  `json:"-" gorm:"type:string;size:128;"` // 为每个用户存一个唯一的jwt key (通用唯一识别码)

	Info map[string]interface{} `json:"info" gorm:"-"`
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) SetPassword(password string) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return
	}
	u.Password = string(b)
}

func (u *User) GetInfo() map[string]interface{} {
	return map[string]interface{}{
		"username":        u.Username,
		"last_login_time": u.LastLoginTime,
	}
}
