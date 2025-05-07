package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"server/core"
	"server/im_group/group_api/internal/config"
	"server/im_group/group_rpc/groups"
	"server/im_group/group_rpc/types/group_rpc"
	"server/im_user/user_rpc/types/user_rpc"
	"server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config   config.Config
	DB       *gorm.DB
	Redis    *redis.Client
	UserRpc  user_rpc.UsersClient
	GroupRpc group_rpc.GroupsClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:   c,
		DB:       mysqlDb,
		Redis:    client,
		UserRpc:  users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		GroupRpc: groups.NewGroups(zrpc.MustNewClient(c.GroupRpc)),
	}
}
