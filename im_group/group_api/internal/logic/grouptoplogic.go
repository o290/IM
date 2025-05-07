package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupTopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupTopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupTopLogic {
	return &GroupTopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupTopLogic) GroupTop(req *types.GroupTopRequest) (resp *types.GroupTopResponse, err error) {
	// 这个群的成员才能调用
	//1.根据群id和用户id查询群信息，能够获取到数据说明是群成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}

	//2.获取群置顶信息，能够获取到说明之前已经置顶，则取消置顶；不能获取到说明之前没有置顶，则置顶
	var userTop group_models.GroupUserTopModel
	err1 := l.svcCtx.DB.Take(&userTop, "group_id = ? and user_id = ?", req.GroupID, req.UserID).Error
	if err1 != nil {
		// 查不到,还没有置顶
		if req.IsTop {
			// 我要置顶
			l.svcCtx.DB.Create(&group_models.GroupUserTopModel{
				GroupID: req.GroupID,
				UserID:  req.UserID,
			})
		}
	} else {
		// 查得到
		if !req.IsTop {
			//取消置顶
			l.svcCtx.DB.Delete(&userTop)
		}
	}
	return
}
