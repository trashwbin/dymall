package dal

import (
	"github.com/trashwbin/dymall/app/ai/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/ai/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
