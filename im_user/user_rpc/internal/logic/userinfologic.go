package logic

import (
	"context"
	"encoding/json"
	"errors"
	"server/im_user/user_models"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *user_rpc.UserInfoRequest) (*user_rpc.UserInfoResponse, error) {
	//1.查找用户信息
	//user用于存储查询结果
	var user user_models.UserModel
	//根据 in.UserId查找id= in.UserId的记录
	//SELECT * FROM users WHERE id = ? LIMIT 1;
	err := l.svcCtx.DB.Preload("UserConfModel").Take(&user, in.UserId).Error
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	//2.将user序列化为json格式的字节切片byteData
	byteData, _ := json.Marshal(user)
	return &user_rpc.UserInfoResponse{Data: byteData}, nil
}
