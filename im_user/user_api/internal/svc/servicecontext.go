package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"net/http"
	"server/core"
	"server/im_chat/chat_rpc/chat"
	"server/im_chat/chat_rpc/types/chat_rpc"
	"server/im_user/user_api/internal/config"
	"server/im_user/user_api/internal/middleware"
	"server/im_user/user_rpc/types/user_rpc"
	"server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	UserRpc         user_rpc.UsersClient
	Redis           *redis.Client
	ChatRpc         chat_rpc.ChatClient
	AdminMiddleware func(next http.HandlerFunc) http.HandlerFunc
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:          c,
		DB:              mysqlDb,
		UserRpc:         users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		Redis:           client,
		ChatRpc:         chat.NewChat(zrpc.MustNewClient(c.ChatRpc)),
		AdminMiddleware: (&middleware.AdminMiddleware{}).Handle,
	}
}
