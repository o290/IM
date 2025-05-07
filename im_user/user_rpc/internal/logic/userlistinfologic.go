package logic

import (
	"context"
	"server/im_user/user_models"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserListInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserListInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserListInfoLogic {
	return &UserListInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserListInfo 根据传入的用户id列表，从数据库中查询这些用户信息
func (l *UserListInfoLogic) UserListInfo(in *user_rpc.UserListInfoRequest) (*user_rpc.UserListInfoResponse, error) {
	//1.查询用户信息
	var userList []user_models.UserModel
	l.svcCtx.DB.Find(&userList, in.UserIdList)
	//2.初始化响应结构体
	resp := new(user_rpc.UserListInfoResponse)
	resp.UserInfo = make(map[uint32]*user_rpc.UserInfo, 0)
	//3.遍历查询结果，封装查询结果到响应体中
	for _, model := range userList {
		resp.UserInfo[uint32(model.ID)] = &user_rpc.UserInfo{
			NickName: model.Nickname,
			Avatar:   model.Avatar,
		}
	}

	return resp, nil
}
