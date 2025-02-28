package mysql

import (
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/user"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo 创建新的 UserRepo 实例
func NewUserRepo() *UserRepo {
	return &UserRepo{
		db: DB, // 假设你有一个 InitDB 函数用于初始化数据库连接
	}
}

// GetUserByUsername 根据用户名获取用户信息
func (repo *UserRepo) GetUserByUsername(username string) (*user.User, error) {
	var u user.User
	if err := repo.db.Where("username = ?", username).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 用户未找到
		}
		return nil, err // 查询失败
	}
	return &u, nil
}

// CreateUser 创建新用户
func (repo *UserRepo) CreateUser(newUser *user.User) error {
	if err := repo.db.Create(newUser).Error; err != nil {
		return err // 插入失败
	}
	return nil
}
