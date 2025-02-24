package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/payment/conf"
	paymentutils "github.com/trashwbin/dymall/app/payment/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth/authservice"
)

var (
	AuthClient authservice.Client
	once       sync.Once
)

func InitClient() {
	once.Do(func() {
		initAuthClient()
	})
}

func initAuthClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	paymentutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	AuthClient, err = authservice.NewClient("auth", opts...)
	paymentutils.MustHandleError(err)
}
