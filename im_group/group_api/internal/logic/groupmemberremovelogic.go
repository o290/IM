package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupMemberRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupMemberRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberRemoveLogic {
	return &GroupMemberRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberRemoveLogic) GroupMemberRemove(req *types.GroupMemberRemoveRequest) (resp *types.GroupMemberRemoveResponse, err error) {
	// 谁能调这个接口 必须得是这个群的成员
	//1.查询调用者的群成员信息
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("违规调用")
	}
	//2.如果退群对象是调用者，先判断角色，不能是群主，其他角色直接删除记录，并往群验证表中添加一条记录
	if req.UserID == req.MemberID {
		// 自己不能是群主 群主不能退群,群主只能解散群
		if member.Role == 1 {
			return nil, errors.New("群主不能退群, 只能解散群聊")
		}
		//把member中的与这个用户的记录删掉就好了
		l.svcCtx.DB.Delete(&member)
		// 给群验证表里面加条记录
		err = l.svcCtx.DB.Create(&group_models.GroupVerifyModel{
			GroupID: member.GroupID,
			UserID:  req.UserID,
			Type:    2, //退群
		}).Error
		return
	}
	//把用户退出群聊
	//3.判断调用者的角色
	if !(member.Role == 1 || member.Role == 2) {
		return nil, errors.New("违规调用")
	}
	//4.查询退群对象的群成员信息
	var member1 group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member1, "group_id = ? and user_id = ?", req.ID, req.MemberID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}

	// 群主可以踢管理员和用户
	// 管理员只能踢用户
	if !(member.Role == 1 && (member1.Role == 2 || member1.Role == 3) || member.Role == 2 && member1.Role == 3) {
		return nil, errors.New("角色错误")
	}
	//5.删除群成员
	err = l.svcCtx.DB.Delete(&member1).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群成员移出失败")
	}
	return
}
