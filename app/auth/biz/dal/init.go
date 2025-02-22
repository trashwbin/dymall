package dal

import (
	"github.com/trashwbin/dymall/app/auth/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
