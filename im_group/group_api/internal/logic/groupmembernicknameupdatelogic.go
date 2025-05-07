package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberNicknameUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberNicknameUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberNicknameUpdateLogic {
	return &GroupMemberNicknameUpdateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberNicknameUpdateLogic) GroupMemberNicknameUpdate(req *types.GroupMemberNicknameUpdateRequest) (resp *types.GroupMemberNicknameUpdateResponse, err error) {
	//1.查询请求发起者是否是群成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("违规调用")
	}
	//2.查询修改昵称的成员是否是群成员
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.ID, req.MemberID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}

	//3.自己修改自己的
	if req.UserID == req.MemberID {
		l.svcCtx.DB.Model(&member).Updates(map[string]any{
			"member_nickname": req.Nickname,
		})
		return
	}

	//4.只有发起者是群主或者是管理员才可以调用 ，群主能修改管理员和普通成员，管理员只能修改普通成员
	if !((member.Role == 1 && (member1.Role == 2 || member1.Role == 3)) || (member.Role == 2 && member1.Role == 3)) {
		return nil, errors.New("用户角色错误")
	}
	l.svcCtx.DB.Model(&member1).Updates(map[string]any{
		"member_nickname": req.Nickname,
	})
	return
}
