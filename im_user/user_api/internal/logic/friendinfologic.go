package logic

import (
	"context"
	"encoding/json"
	"errors"
	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"
	"server/im_user/user_models"
	"server/im_user/user_rpc/types/user_rpc"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendInfoLogic {
	return &FriendInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendInfoLogic) FriendInfo(req *types.FriendInfoRequest) (resp *types.FriendInfoResponse, err error) {
	var friend user_models.FriendModel
	//1.判断是否是自己的好友
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("他不是你的好友")
	}
	//2.获取好友的用户信息
	res, err := l.svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
		UserId: uint32(req.FriendID),
	})
	if err != nil {
		return nil, errors.New(err.Error())
	}
	//3.判断好友是否在线
	onlineMap := l.svcCtx.Redis.HGetAll(l.ctx, "online").Val()
	//将req.FriendID转化为字符串类型
	//不能直接使用string强制转换，因为string强制转换是用于字节切片转成字符串类型的
	//整数转换成字符串类型是需要使用 strconv.Itoa方法
	key := strconv.Itoa(int(req.FriendID))
	_, ok := onlineMap[key]
	//4.将数据反序列化
	var friendUser user_models.UserModel
	json.Unmarshal(res.Data, &friendUser)

	//5.构造响应
	response := types.FriendInfoResponse{
		UserID:   friendUser.ID,
		NickName: friendUser.Nickname,
		Abstract: friendUser.Abstract,
		Avatar:   friendUser.Avatar,
		Notice:   friend.GetUserNotice(req.UserID),
		IsOline:  ok,
	}
	return &response, nil
}
