package logic

import (
	"context"

	"server/im_group/group_rpc/internal/svc"
	"server/im_group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsInGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsInGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsInGroupLogic {
	return &IsInGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsInGroupLogic) IsInGroup(in *group_rpc.IsInGroupRequest) (*group_rpc.IsInGroupResponse, error) {
	// todo: add your logic here and delete this line

	return &group_rpc.IsInGroupResponse{}, nil
}
