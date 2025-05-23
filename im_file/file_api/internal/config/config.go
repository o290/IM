package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Etcd      string
	FileSize  float64
	BlackList []string
	WhiteList []string
	UploadDir string //上传文件保存的目录
	UserRpc   zrpc.RpcClientConf
	Mysql     struct {
		DataSource string
	}
}
