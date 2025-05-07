package logic

import (
	"context"
	"server/im_user/user_models"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsFriendLogic {
	return &IsFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsFriendLogic) IsFriend(in *user_rpc.IsFriendRequest) (res *user_rpc.IsFriendResponse, err error) {
	res = new(user_rpc.IsFriendResponse)
	var friend user_models.FriendModel
	if friend.IsFriend(l.svcCtx.DB, uint(in.User1), uint(in.User2)) {
		res.IsFriend = true
		return
	}

	return
}
