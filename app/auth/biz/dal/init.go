package dal

import (
	"github.com/trashwbin/dymall/app/auth/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/auth/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
