package svc

import (
    "github.com/shorturl/short-url/zero_remake/shorturl_api/internal/config"
    "github.com/shorturl/short-url/zero_remake/shorturl_rpc/shorturlclient"
    "github.com/shorturl/short-url/zero_remake/user-rpc/userservice"
    "github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	ShortUrlRpc shorturlclient.ShortUrl
	UserRpc     userservice.UserService
	LogPath     string
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		ShortUrlRpc: shorturlclient.NewShortUrl(zrpc.MustNewClient(c.ShortUrlRpc)),
		UserRpc:     userservice.NewUserService(zrpc.MustNewClient(c.UserRpc)),
		LogPath:     c.Log.Path,
	}
}
