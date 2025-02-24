package dal

import (
	"github.com/trashwbin/dymall/app/product/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/product/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
