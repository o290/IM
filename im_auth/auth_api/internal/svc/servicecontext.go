package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"server/core"
	"server/im_auth/auth_api/internal/config"
	"server/im_user/user_rpc/types/user_rpc"
	"server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB //注入
	Redis   *redis.Client
	UserRpc user_rpc.UsersClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	//连接mysql
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	//mysqlDb.AutoMigrate(&auth_models.UserConfModel{},&auth_models.FriendModel{},&auth_models.UserModel{},&auth_models.FriendVerifyModel{})
	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   client,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
	}
}
