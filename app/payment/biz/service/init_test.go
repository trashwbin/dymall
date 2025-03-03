package service

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/trashwbin/dymall/app/payment/biz/dal"
	"github.com/trashwbin/dymall/app/payment/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/payment/biz/model"
)

var (
	testPayment *model.Payment
)

func TestMain(m *testing.M) {
	// 加载环境变量
	_ = godotenv.Load()

	// 初始化数据库和缓存
	dal.Init()

	// 清理可能存在的测试数据
	if err := cleanupTestData(); err != nil {
		panic(err)
	}

	// 创建测试数据
	if err := setupTestData(); err != nil {
		panic(err)
	}

	// 运行测试
	code := m.Run()

	// 清理测试数据
	if err := cleanupTestData(); err != nil {
		panic(err)
	}

	os.Exit(code)
}

// 清理测试数据
func cleanupTestData() error {
	repo := mysql.NewPaymentRepo()
	return repo.DeletePaymentForTest("test_payment_id")
}

// 创建测试数据
func setupTestData() error {
	repo := mysql.NewPaymentRepo()

	// 创建测试支付单
	now := time.Now()
	payment := &model.Payment{
		PaymentID: "test_payment_id",
		OrderID:   "test_order_id",
		UserID:    1001,
		Amount:    99.9,
		Currency:  "CNY",
		Status:    model.PayStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
		ExpireAt:  now.Add(24 * time.Hour),
	}

	var err error
	testPayment, err = repo.CreatePayment(payment)
	if err != nil {
		return err
	}

	return nil
}
