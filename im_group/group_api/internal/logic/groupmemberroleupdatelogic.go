package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberRoleUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberRoleUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberRoleUpdateLogic {
	return &GroupMemberRoleUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberRoleUpdateLogic) GroupMemberRoleUpdate(req *types.GroupMemberRoleUpdateRequest) (resp *types.GroupMemberRoleUpdateResponse, err error) {
	//1.查询调用者的群成员信息
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("违规调用")
	}
	//2.调用者需是管理员
	if member.Role != 1 {
		return nil, errors.New("权限错误")
	}
	//3.查询修改对象的群成员信息
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.ID, req.MemberID).Error
	if err != nil {
		return nil, errors.New("用户还不是群成员呢")
	}
	//4.修改对象必须是管理员或普通成员
	if !(req.Role == 2 || req.Role == 3) {
		return nil, errors.New("用户角色设置错误")
	}
	if member1.Role == req.Role {
		return
	}
	//5.修改角色
	l.svcCtx.DB.Model(&member1).Update("role", req.Role)
	return
}
