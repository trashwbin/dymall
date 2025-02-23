package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/cart/conf"
	cartutils "github.com/trashwbin/dymall/app/cart/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth/authservice"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product/productcatalogservice"
)

var (
	ProductClient productcatalogservice.Client
	AuthClient    authservice.Client
	once          sync.Once
)

func InitClient() {
	once.Do(func() {
		initProductClient()
		initAuthClient()
	})
}

func initProductClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	cartutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	cartutils.MustHandleError(err)
}

func initAuthClient() {
	var opts []client.Option
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	cartutils.MustHandleError(err)
	opts = append(opts, client.WithResolver(r))
	AuthClient, err = authservice.NewClient("auth", opts...)
	cartutils.MustHandleError(err)
}
