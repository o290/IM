package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"
	"server/utils/set"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryDeleteLogic {
	return &GroupHistoryDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupHistoryDeleteLogic) GroupHistoryDelete(req *types.GroupHistoryDeleteRequest) (resp *types.GroupHistoryDeleteListResponse, err error) {
	//1.查询该用户是否是群的成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}
	//2.查询我删除了哪些聊天记录
	var msgIDList []uint
	l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).
		Where("group_id = ? and user_id = ?", req.ID, req.UserID).
		Select("msg_id").Scan(&msgIDList)

	//3.和我要删除的聊天记录求差集，计算出真正要删除的记录
	addMsgIDList := set.Difference(req.MsgIDList, msgIDList)
	logx.Infof("删除聊天记录的id列表 %v", addMsgIDList)
	if len(addMsgIDList) == 0 {
		return
	}

	//4.获取要删除的聊天记录消息，判断是否存在，用户传过来的消息id不一定存在
	var msgIDFindList []uint
	l.svcCtx.DB.Model(group_models.GroupMsgModel{}).
		Where("id in ?", addMsgIDList).
		Select("id").Scan(&msgIDFindList)
	if len(msgIDFindList) != len(addMsgIDList) {
		return nil, errors.New("消息一致性错误")
	}
	//4.删除聊天记录（往聊天记录表中添加）
	var list []group_models.GroupUserMsgDeleteModel
	for _, i2 := range addMsgIDList {
		list = append(list, group_models.GroupUserMsgDeleteModel{
			MsgID:   i2,
			UserID:  req.UserID,
			GroupID: req.ID,
		})
	}
	err = l.svcCtx.DB.Create(&list).Error
	if err != nil {
		return
	}
	return
}
