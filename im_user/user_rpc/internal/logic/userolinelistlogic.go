package logic

import (
	"context"

	"server/im_user/user_rpc/internal/svc"
	"server/im_user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOlineListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOlineListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOlineListLogic {
	return &UserOlineListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserOlineListLogic) UserOlineList(in *user_rpc.UserOlineRequest) (*user_rpc.UserOlineResponse, error) {
	// todo: add your logic here and delete this line

	return &user_rpc.UserOlineResponse{}, nil
}
