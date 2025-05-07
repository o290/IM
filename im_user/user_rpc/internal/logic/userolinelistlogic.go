package logic

import (
	"context"
	"strconv"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOlineListLogic {
	return &UserOlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserOlineListLogic) UserOlineList(in *user_rpc.UserOlineListRequest) (resp *user_rpc.UserOlineListResponse, err error) {
	resp = new(user_rpc.UserOlineListResponse)
	//1.获取在线用户列表
	onlineMap := l.svcCtx.Redis.HGetAll(context.Background(), "online").Val()
	//2.遍历每一个键值对，存储在线用户id
	for key, _ := range onlineMap {
		//将类型为字符串的key转换成int类型
		val, err1 := strconv.Atoi(key)
		if err1 != nil {
			logx.Error(err1)
			continue
		}
		//resp.UserIdList用于存储在线用户的id
		resp.UserIdList = append(resp.UserIdList, uint32(val))
	}

	return
}
