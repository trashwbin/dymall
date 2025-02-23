package dal

import (
	"github.com/trashwbin/dymall/app/cart/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
