// #file:D:\Code\Work\dymall\app\user\biz\dal\mysql\init.go
package mysql

import (
	"fmt"
	"os"
	"time"

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
	// 从环境变量中读取MySQL配置
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")

	// 构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/user?charset=utf8mb4&parseTime=True&loc=Local", user, password, host)

	// 初始化GORM连接
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt:            true,                                                         // 开启预编译语句
		SkipDefaultTransaction: true,                                                         // 跳过默认事务
		Logger:                 logger.Default.LogMode(logger.Info),                          // 设置日志级别
		NowFunc:                func() time.Time { return time.Now().Truncate(time.Second) }, // 设置时间精度为秒
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
