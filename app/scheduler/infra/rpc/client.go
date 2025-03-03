package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/scheduler/conf"
	schedulerutils "github.com/trashwbin/dymall/app/scheduler/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth/authservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment/paymentservice"
)

var (
	AuthClient    authservice.Client
	OrderClient   orderservice.Client
	PaymentClient paymentservice.Client
	once          sync.Once
)

func InitClient() {
	once.Do(func() {
		initAuthClient()
		initOrderClient()
		initPaymentClient()
	})
}

func initAuthClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	schedulerutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	AuthClient, err = authservice.NewClient("auth", opts...)
	schedulerutils.MustHandleError(err)
}

func initOrderClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	schedulerutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	OrderClient, err = orderservice.NewClient("order", opts...)
	schedulerutils.MustHandleError(err)
}

func initPaymentClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	schedulerutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	PaymentClient, err = paymentservice.NewClient("payment", opts...)
	schedulerutils.MustHandleError(err)
}
