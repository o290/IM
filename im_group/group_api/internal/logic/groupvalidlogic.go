package logic

import (
	"context"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidLogic {
	return &GroupValidLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidLogic) GroupValid(req *types.GroupValidRequest) (resp *types.GroupValidResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
