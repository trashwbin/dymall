package service

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/cart/biz/dal"
	"github.com/trashwbin/dymall/app/cart/infra/rpc"
)

func TestMain(m *testing.M) {
	// 替换为mock客户端
	rpc.ProductClient = &MockProductClient{}

	// 加载环境变量
	_ = godotenv.Load()

	// 初始化数据库
	dal.Init()

	// 运行测试
	code := m.Run()
	os.Exit(code)
}
