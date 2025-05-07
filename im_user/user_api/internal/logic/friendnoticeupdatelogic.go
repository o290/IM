package logic

import (
	"context"
	"errors"
	"server/im_user/user_models"

	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendNoticeUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFriendNoticeUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendNoticeUpdateLogic {
	return &FriendNoticeUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendNoticeUpdateLogic) FriendNoticeUpdate(req *types.FriendNoticeUpdateRequest) (resp *types.FriendNoticeUpdateResponse, err error) {
	var friend user_models.FriendModel
	//1.判断是否是好友
	if !friend.IsFriend(l.svcCtx.DB, req.UserID, req.FriendID) {
		return nil, errors.New("他不是你的好友")
	}

	//2.判断备注方是谁
	if friend.SendUserID == req.UserID {
		//我是发起方
		//SendUserNotice指的是发送发对接收方的备注
		if friend.SendUserNotice == req.Notice {
			return
		}
		l.svcCtx.DB.Model(&friend).Update("send_user_notice", req.Notice)
	}
	if friend.RevUserID == req.UserID {
		if friend.RevUserNotice == req.Notice {
			return
		}
		l.svcCtx.DB.Model(&friend).Update("rev_user_notice", req.Notice)
	}
	return
}
