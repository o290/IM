package svc

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"server/core"
	"server/im_user/user_rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB //注入
	Redis  *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.UserMysql.DataSource)
	client := core.InitRedis(c.RedisConf.Addr, c.RedisConf.Pwd, c.RedisConf.DB)
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
		Redis:  client,
	}
}
