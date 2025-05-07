package redis_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"server/common/models/ctype"
	"server/im_user/user_rpc/types/user_rpc"
	"time"
)

// GetUserBaseInfo 获取用户信息
// 从 Redis 缓存中获取用户的基本信息。如果缓存中不存在该用户信息，
// 则通过 RPC 调用从用户服务中获取，并将获取到的信息存入 Redis 缓存，以便后续使用
func GetUserBaseInfo(client *redis.Client, userRpc user_rpc.UsersClient, userID uint) (userInfo ctype.UserInfo, err error) {
	key := fmt.Sprintf("fim_server_user_%d", userID)
	str, err := client.Get(context.Background(), key).Result()
	if err != nil {
		// 没找到
		fmt.Println("2255552")
		userBaseResponse, err1 := userRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(userID),
		})
		if err1 != nil {
			err = err1
			return
		}
		err = nil
		userInfo.ID = userID
		userInfo.Avatar = userBaseResponse.Avatar
		userInfo.NickName = userBaseResponse.NickName

		byteData, _ := json.Marshal(userInfo)
		// 设置进缓存
		client.Set(context.Background(), key, string(byteData), time.Hour) //1小时过期
		return
	}

	err = json.Unmarshal([]byte(str), &userInfo)
	if err != nil {
		return
	}
	return
}
