package main

import (
	"example.com/shorturl/short-url/zero_remake/middleware"
	"flag"
	"fmt"

	"example.com/shorturl/short-url/zero_remake/user-rpc/internal/config"
	"example.com/shorturl/short-url/zero_remake/user-rpc/internal/server"
	"example.com/shorturl/short-url/zero_remake/user-rpc/internal/svc"
	"example.com/shorturl/short-url/zero_remake/user-rpc/types/User"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/userrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		User.RegisterUserServiceServer(grpcServer, server.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	s.AddUnaryInterceptors(middleware.UnaryServerLogger(c.Log.Path))
	s.AddStreamInterceptors(middleware.StreamServerLogger(c.Log.Path))

	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
