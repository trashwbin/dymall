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
	var userDO UserDO
	if err := repo.db.Where("username = ?", username).First(&userDO).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // 用户未找到
		}
		return nil, err // 查询失败
	}
	// 直接转换为 user.User
	return &user.User{
		Id:       int64(userDO.ID),
		Username: userDO.Username,
		Password: userDO.Password,
		Email:    userDO.Email,
		Gender:   userDO.Gender,
		Age:      int32(userDO.Age),
		Address:  userDO.Address,
	}, nil
}

// CreateUser 创建新用户
func (repo *UserRepo) CreateUser(newUser *user.User) error {
	userDO := &UserDO{
		Username: newUser.Username,
		Password: newUser.Password,
		Email:    newUser.Email,
		Gender:   newUser.Gender,
		Age:      int(newUser.Age),
		Address:  newUser.Address,
	}

	if err := repo.db.Create(userDO).Error; err != nil {
		return err // 插入失败
	}
	newUser.Id = int64(userDO.ID)
	return nil
}
