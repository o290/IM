package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	UserMysql struct {
		DataSource string
	}
	RedisConf struct {
		Addr string
		Pwd  string
		DB   int
	}
}
