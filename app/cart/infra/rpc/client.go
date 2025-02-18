package rpc

import (
	"sync"

	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/cart/conf"
	cartutils "github.com/trashwbin/dymall/app/cart/utils"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/product/productcatalogservice"
)

var (
	ProductClient productcatalogservice.Client
	once          sync.Once
)

func InitClient() {
	once.Do(func() {
		initProductClient()
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
