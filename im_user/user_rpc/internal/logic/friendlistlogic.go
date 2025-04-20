package logic

import (
	"context"
	"server/common/list_query"
	"server/common/models"
	"server/im_user/user_models"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendListLogic {
	return &FriendListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FriendListLogic) FriendList(in *user_rpc.FriendListRequest) (*user_rpc.FriendListResponse, error) {
	friends, _, _ := list_query.ListQuery(l.svcCtx.DB, user_models.FriendModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Limit: -1, //查全部
		},
		Preload: []string{"SendUserModel", "RevUserModel"},
	})

	var list []*user_rpc.FriendInfo
	for _, friend := range friends {
		info := user_rpc.FriendInfo{}
		if friend.SendUserID == uint(in.User) {
			//我是发起方,返回接收方
			info = user_rpc.FriendInfo{
				UserId:   uint32(friend.RevUserID),
				NickName: friend.RevUserModel.Nickname,
				Avatar:   friend.RevUserModel.Avatar,
			}
		}
		if friend.RevUserID == uint(in.User) {
			//我是发起方
			info = user_rpc.FriendInfo{
				UserId:   uint32(friend.SendUserID),
				NickName: friend.SendUserModel.Nickname,
				Avatar:   friend.SendUserModel.Avatar,
			}
		}
		list = append(list, &info)
	}

	return &user_rpc.FriendListResponse{FriendList: list}, nil
}
