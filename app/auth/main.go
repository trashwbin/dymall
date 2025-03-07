package main

import (
	"github.com/joho/godotenv"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/trashwbin/dymall/app/auth/biz/dal"
	"github.com/trashwbin/dymall/app/auth/biz/dal/mysql"
	"github.com/trashwbin/dymall/app/auth/biz/middleware"
	"github.com/trashwbin/dymall/app/auth/biz/service"
	"github.com/trashwbin/dymall/app/auth/conf"
	"github.com/trashwbin/dymall/rpc_gen/kitex_gen/auth/authservice"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	_ = godotenv.Load()
	dal.Init()
	opts := kitexInit()

	// 初始化授权服务
	authSvc, err := service.NewAuthorizationService(mysql.DB)
	if err != nil {
		klog.Fatal(err)
	}

	// 添加认证中间件
	opts = append(opts, server.WithMiddleware(middleware.AuthMiddleware(authSvc)))

	// 创建服务实现
	impl := NewAuthServiceImpl(authSvc)

	svr := authservice.NewServer(impl, opts...)
	err = svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))
	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))
	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}
