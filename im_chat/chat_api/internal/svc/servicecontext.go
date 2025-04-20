package svc

import (
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
	"server/core"
	"server/im_chat/chat_api/internal/config"
	"server/im_chat/chat_rpc/chat"
	"server/im_chat/chat_rpc/types/chat_rpc"
	"server/im_file/file_rpc/files"
	"server/im_file/file_rpc/types/file_rpc"
	"server/im_user/user_rpc/types/user_rpc"
	"server/im_user/user_rpc/users"
)

type ServiceContext struct {
	Config  config.Config
	DB      *gorm.DB
	Redis   *redis.Client
	UserRpc user_rpc.UsersClient
	FileRpc file_rpc.FilesClient
	ChatRpc chat_rpc.ChatClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	client := core.InitRedis(c.Redis.Addr, c.Redis.Pwd, c.Redis.DB)
	return &ServiceContext{
		Config:  c,
		DB:      mysqlDb,
		Redis:   client,
		UserRpc: users.NewUsers(zrpc.MustNewClient(c.UserRpc)),
		FileRpc: files.NewFiles(zrpc.MustNewClient(c.FileRpc)),
		ChatRpc: chat.NewChat(zrpc.MustNewClient(c.ChatRpc)),
	}
}
