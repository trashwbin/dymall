// #file:D:\Code\Work\dymall\app\user\biz\model\user.go
package model

import (
	"time"
)

// Gender 用户性别枚举类型
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// User 领域模型
type User struct {
	ID        uint      `json:"id"`         // 用户ID
	Username  string    `json:"username"`   // 用户名
	Password  string    `json:"password"`   // 密码
	Email     string    `json:"email"`      // 邮箱
	Gender    Gender    `json:"gender"`     // 性别
	Age       int       `json:"age"`        // 年龄
	Address   string    `json:"address"`    // 地址
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}
