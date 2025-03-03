package dal

import (
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/scheduler/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
