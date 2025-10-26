package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DbUser string
		DbPort string
		DbPass string
		DbName string
		DbHost string
	}
	log struct {
		path string
	}
}
