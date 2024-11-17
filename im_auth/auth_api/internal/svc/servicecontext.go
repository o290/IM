package svc

import (
	"gorm.io/gorm"
	"server/core"
	"server/im_auth/auth_api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB //注入
}

func NewServiceContext(c config.Config) *ServiceContext {
	//连接mysql
	mysqlDb := core.InitMysql(c.Mysql.DataSource)
	//mysqlDb.AutoMigrate(&auth_models.UserConfModel{},&auth_models.FriendModel{},&auth_models.UserModel{},&auth_models.FriendVerifyModel{})
	return &ServiceContext{
		Config: c,
		DB:     mysqlDb,
	}
}
