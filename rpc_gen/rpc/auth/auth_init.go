package auth

import (
	"sync"

	"github.com/cloudwego/kitex/client"
)

var (
	// todo edit custom config
	defaultClient     RPCClient
	defaultDstService = "auth"
	defaultClientOpts = []client.Option{
		client.WithHostPorts("192.168.148.120:8880"),
	}
	once sync.Once
)

func init() {
	DefaultClient()
}

func DefaultClient() RPCClient {
	once.Do(func() {
		defaultClient = newClient(defaultDstService, defaultClientOpts...)
	})
	return defaultClient
}

func newClient(dstService string, opts ...client.Option) RPCClient {
	c, err := NewRPCClient(dstService, opts...)
	if err != nil {
		panic("failed to init client: " + err.Error())
	}
	return c
}

func InitClient(dstService string, opts ...client.Option) {
	defaultClient = newClient(dstService, opts...)
}
