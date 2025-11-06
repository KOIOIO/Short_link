package svc

import (
    "github.com/shorturl/short-url/zero_remake/common/init_gorm"
    "github.com/shorturl/short-url/zero_remake/user-rpc/internal/config"
    "fmt"
    "gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Mysql.DbUser,
		c.Mysql.DbPass,
		c.Mysql.DbHost,
		c.Mysql.DbPort,
		c.Mysql.DbName,
	)
	db := init_gorm.Init_gorm(dns)

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
