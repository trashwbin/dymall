package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/product/conf"
	productutils "github.com/trashwbin/dymall/app/product/utils"
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
	productutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	AuthClient, err = authservice.NewClient("auth", opts...)
	productutils.MustHandleError(err)
}
