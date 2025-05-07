package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupRemoveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupRemoveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupRemoveLogic {
	return &GroupRemoveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupRemoveLogic) GroupRemove(req *types.GroupRemoveRequest) (resp *types.GroupRemoveResponse, err error) {
	//1.查询该用户是否是群成员
	var groupMember group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&groupMember, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("群不存在或用户不是群成员")
	}
	//2.判断是否是群主
	if groupMember.Role != 1 {
		return nil, errors.New("只能群主才能解散该群")
	}

	//3.删除群信息
	var msgList []group_models.GroupMsgModel
	l.svcCtx.DB.Find(&msgList, "group_id = ?", req.ID).Delete(&msgList)
	//4.删除群成员
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ?", req.ID).Delete(&memberList)
	//5.删除群验证消息
	var vList []group_models.GroupVerifyModel
	l.svcCtx.DB.Find(&vList, "group_id = ?", req.ID).Delete(&vList)
	//6.删除群
	var group group_models.GroupModel
	l.svcCtx.DB.Take(&group, req.ID).Delete(&group)

	logx.Infof("删除群:", group.Title)
	logx.Infof("关联群成员数:", len(memberList))
	logx.Infof("关联群消息数:", len(msgList))
	logx.Infof("关联群验证消息数:", len(vList))
	return
}
