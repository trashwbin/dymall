package dal

import (
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
