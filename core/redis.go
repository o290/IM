package core

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

var rdb *redis.Client

// Background返回一个非空context，它永远不会取消，没有值
// 也没有期限，通常在main函数，初始化和测试时使用，并用作
// 传入请求的顶级上下文
var ctx = context.Background()

func InitRedis(addr, pwd string, db int) (client *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,  //默认db为0，指的是连接到哪一个数据库
		PoolSize: 100, //连接池大小
	})
	_, cancel := context.WithTimeout(ctx, 500*time.Microsecond)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
		return
	}
	return rdb
}
