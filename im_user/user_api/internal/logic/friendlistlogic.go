package logic

import (
	"context"
	"server/common/list_query"
	"server/common/models"
	"server/im_user/user_models"
	"strconv"

	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FriendList 获取好友列表信息
func (l *FriendListLogic) FriendList(req *types.FriendListRequest) (resp *types.FriendListResponse, err error) {
	//1.查询该用户的好友列表信息和好友总数
	friends, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		//分页查询
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		//预加载
		Preload: []string{"SendUserModel", "RevUserModel"},
		Where:   l.svcCtx.DB.Where("send_user_id = ? or rev_user_id = ?", req.UserID, req.UserID),
	})
	//2.从redis中获取在线用户信息
	onlineMap := l.svcCtx.Redis.HGetAll(l.ctx, "online").Val()
	var onlineUserMap = map[uint]bool{}
	for key, _ := range onlineMap {
		//将在线用户id转换成uint类型
		val, err1 := strconv.Atoi(key)
		if err1 != nil {
			logx.Error(err1)
			continue
		}
		//存储到map中
		onlineUserMap[uint(val)] = true
	}

	//3.封装好友信息
	var list []types.FriendInfoResponse
	for _, friend := range friends {
		info := types.FriendInfoResponse{}
		if friend.SendUserID == req.UserID {
			//我是发起方,返回接收方
			info = types.FriendInfoResponse{
				UserID:   friend.RevUserID,
				NickName: friend.RevUserModel.Nickname,
				Abstract: friend.RevUserModel.Abstract,
				Avatar:   friend.RevUserModel.Avatar,
				Notice:   friend.SendUserNotice, //我是发起方，我给接收方备注
				IsOline:  onlineUserMap[friend.RevUserID],
			}
		}
		if friend.RevUserID == req.UserID {
			//我是发起方
			info = types.FriendInfoResponse{
				UserID:   friend.SendUserID,
				NickName: friend.SendUserModel.Nickname,
				Abstract: friend.SendUserModel.Abstract,
				Avatar:   friend.SendUserModel.Avatar,
				Notice:   friend.RevUserNotice,
				IsOline:  onlineUserMap[friend.SendUserID],
			}
		}
		list = append(list, info)
	}
	return &types.FriendListResponse{List: list, Count: int(count)}, nil
}
