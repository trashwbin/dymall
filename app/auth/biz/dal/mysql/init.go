package mysql

import (
	"fmt"
	"os"

	"github.com/trashwbin/dymall/app/auth/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {

	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))

	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err)
	}

	// 自动迁移 Casbin 规则表
	err = DB.AutoMigrate(&CasbinRule{})
	if err != nil {
		panic(err)
	}
}
