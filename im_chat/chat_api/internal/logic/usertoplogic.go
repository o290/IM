package logic

import (
	"context"
	"errors"
	"server/im_chat/chat_models"
	"server/im_user/user_rpc/types/user_rpc"

	"server/im_chat/chat_api/internal/svc"
	"server/im_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserTopLogic {
	return &UserTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserTop 用户置顶
func (l *UserTopLogic) UserTop(req *types.UserTopRequest) (resp *types.UserTopResponse, err error) {
	//1.判断是否是好友
	if req.UserID != req.FriendID {
		//判断是否是好友
		res, err := l.svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
			User2: uint32(req.UserID),
			User1: uint32(req.FriendID),
		})
		if err != nil {
			return nil, err
		}
		if !res.IsFriend {
			return nil, errors.New("你们还不是好友")
		}
	}

	var topUser chat_models.TopUserModel
	err1 := l.svcCtx.DB.Take(&topUser, "user_id=? and top_user_id=?", req.UserID, req.FriendID).Error
	if err1 != nil {
		//没有置顶
		l.svcCtx.DB.Create(&chat_models.TopUserModel{
			UserID:    req.UserID,
			TopUserID: req.FriendID,
		})
		return
	}
	//已经有置顶了
	l.svcCtx.DB.Delete(&topUser)
	return
}
