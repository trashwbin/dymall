// #file:D:\Code\Work\dymall\app\user\biz\dal\mysql\init.go
package mysql

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	err error
)

// Init 初始化MySQL连接
func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/user?charset=utf8mb4&parseTime=True&loc=Local", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(fmt.Errorf("连接数据库失败: %w", err))
	}

	// 在非生产环境下自动迁移
	if os.Getenv("GO_ENV") != "online" {
		err = DB.AutoMigrate(&UserDO{})
		if err != nil {
			panic(fmt.Errorf("数据库迁移失败: %w", err))
		}
	}

	// 获取底层数据库实例
	sqlDB, err := DB.DB()
	if err != nil {
		panic(fmt.Errorf("获取数据库实例失败: %w", err))
	}

	// 设置连接池
	sqlDB.SetMaxIdleConns(10)  // 最大空闲连接数
	sqlDB.SetMaxOpenConns(100) // 最大打开连接数
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
