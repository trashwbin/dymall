package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/order/conf"
	orderutils "github.com/trashwbin/dymall/app/order/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth/authservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/scheduler/schedulerservice"
)

var (
	AuthClient      authservice.Client
	CartClient      cartservice.Client
	ProductClient   productcatalogservice.Client
	SchedulerClient schedulerservice.Client
	PaymentClient   paymentservice.Client
	once            sync.Once
)

func InitClient() {
	once.Do(func() {
		initAuthClient()
		initCartClient()
		initProductClient()
		initSchedulerClient()
		initPaymentClient()
	})
}

func initAuthClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Errorf("create consul resolver failed: %v", err)
		orderutils.MustHandleError(err)
	}
	opts = append(opts, client.WithResolver(r))
	AuthClient, err = authservice.NewClient("auth", opts...)
	if err != nil {
		klog.Errorf("create auth client failed: %v", err)
		orderutils.MustHandleError(err)
	}
}

func initCartClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Errorf("create consul resolver failed: %v", err)
		orderutils.MustHandleError(err)
	}
	opts = append(opts, client.WithResolver(r))
	CartClient, err = cartservice.NewClient("cart", opts...)
	if err != nil {
		klog.Errorf("create cart client failed: %v", err)
		orderutils.MustHandleError(err)
	}
}

func initProductClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Errorf("create consul resolver failed: %v", err)
		orderutils.MustHandleError(err)
	}
	opts = append(opts, client.WithResolver(r))
	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	if err != nil {
		klog.Errorf("create product client failed: %v", err)
		orderutils.MustHandleError(err)
	}
}

func initSchedulerClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Errorf("create consul resolver failed: %v", err)
		orderutils.MustHandleError(err)
	}
	opts = append(opts, client.WithResolver(r))
	SchedulerClient, err = schedulerservice.NewClient("scheduler", opts...)
	if err != nil {
		klog.Errorf("create scheduler client failed: %v", err)
		orderutils.MustHandleError(err)
	}
}

func initPaymentClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Errorf("create consul resolver failed: %v", err)
		orderutils.MustHandleError(err)
	}
	opts = append(opts, client.WithResolver(r))
	PaymentClient, err = paymentservice.NewClient("payment", opts...)
	if err != nil {
		klog.Errorf("create payment client failed: %v", err)
		orderutils.MustHandleError(err)
	}
}
