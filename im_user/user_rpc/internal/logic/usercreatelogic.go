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
	//1.查看用户的openid是否存在，存在则不需要创建
	var user user_models.UserModel
	err := l.svcCtx.DB.Take(&user, "open_id=?", in.OpenId).Error
	if err != nil {
		return nil, errors.New("该用户已存在")
	}
	//2.构造用户信息结构体
	user = user_models.UserModel{
		Nickname:       in.NickName,
		Avatar:         in.Avatar,
		Role:           int8(in.Role),
		OpenID:         in.OpenId,
		RegisterSource: in.RegisterSource,
	}
	//3.创建用户
	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("创建用户失败")
	}
	//4.创建用户配置
	l.svcCtx.DB.Create(&user_models.UserConfModel{
		UserID:        user.ID,
		RecallMessage: nil,
		FriendOnline:  false,
		Sound:         true,
		SecureLink:    false,
		SavePwd:       false,
		SearchUser:    2,
		Verification:  2,
		Online:        true,
	})
	return &user_rpc.UserCreateResponse{UserId: int32(user.ID)}, nil
}
