package dal

import (
	"github.com/trashwbin/dymall/app/checkout/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
