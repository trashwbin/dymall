// #file:D:\Code\Work\dymall\app\user\biz\dal\mysql\user_do.go
package mysql

import (
	"time"

	"github.com/trashwbin/dymall/app/user/biz/module"
	"gorm.io/gorm"
)

// UserDO 用户表模型
type UserDO struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`                                        // 用户ID
	Username  string         `gorm:"unique;not null;comment:用户名"`                                     // 用户名
	Password  string         `gorm:"not null;comment:密码"`                                             // 密码
	Email     string         `gorm:"comment:邮箱"`                                                      // 邮箱
	Gender    string         `gorm:"type:enum('male', 'female', 'other');default:'other';comment:性别"` // 性别
	Age       int            `gorm:"comment:年龄"`                                                      // 年龄
	Address   string         `gorm:"comment:地址"`                                                      // 地址
	CreatedAt time.Time      `gorm:"not null;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"not null;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间"` // 删除时间
}

// TableName 设置表名
func (UserDO) TableName() string {
	return "user"
}

// ToModel 转换为领域模型
func (u *UserDO) ToModel() *model.User {
	return &model.User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		Gender:    model.Gender(u.Gender),
		Age:       u.Age,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// FromModel 从领域模型转换
func (u *UserDO) FromModel(m *model.User) {
	u.ID = m.ID
	u.Username = m.Username
	u.Password = m.Password
	u.Email = m.Email
	u.Gender = string(m.Gender)
	u.Age = m.Age
	u.Address = m.Address
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
}
