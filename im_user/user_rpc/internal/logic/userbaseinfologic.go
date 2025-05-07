package logic

import (
	"context"
	"errors"
	"server/im_user/user_models"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBaseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBaseInfoLogic {
	return &UserBaseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserBaseInfoLogic) UserBaseInfo(in *user_rpc.UserBaseInfoRequest) (*user_rpc.UserBaseInfoResponse, error) {
	//1.查找该用户的用户信息
	var user user_models.UserModel
	err := l.svcCtx.DB.Take(&user, in.UserId).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	//2.返回响应
	return &user_rpc.UserBaseInfoResponse{
		UserId:   in.UserId,
		Avatar:   user.Nickname,
		NickName: user.Avatar,
	}, nil
}
