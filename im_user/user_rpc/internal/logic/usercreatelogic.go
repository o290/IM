package logic

import (
	"context"
	"errors"
	"server/im_user/user_models"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserCreateLogic) UserCreate(in *user_rpc.UserCreateRequest) (*user_rpc.UserCreateResponse, error) {
	//查看用户的openid是否存在
	var user user_models.UserModel
	err := l.svcCtx.DB.Take(&user, "open_id=?", in.OpenId).Error
	if err != nil {
		return nil, errors.New("该用户已存在")
	}
	user = user_models.UserModel{
		Nickname:       in.NickName,
		Avatar:         in.Avatar,
		Role:           int8(in.Role),
		OpenID:         in.OpenId,
		RegisterSource: in.RegisterSource,
	}
	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("创建用户失败")
	}
	return &user_rpc.UserCreateResponse{UserId: int32(user.ID)}, nil
}
