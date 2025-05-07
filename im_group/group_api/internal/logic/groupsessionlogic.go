package logic

import (
	"context"
	"fmt"
	"server/common/list_query"
	"server/common/models"
	"server/im_group/group_models"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSessionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSessionLogic {
	return &GroupSessionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type SessionData struct {
	GroupID       uint   `gorm:"column:group_id"`
	NewMsgDate    string `gorm:"column: newMsgDate"`
	NewMsgPreview string `gorm:"column:newMsgPreview"`
	IsTop         bool   `gorm:"column:isTop"`
}

func (l *GroupSessionLogic) GroupSession(req *types.GroupSessionRequest) (resp *types.GroupSessionListResponse, err error) {
	// 1.查询出我加入的群聊id列表
	var userGroupIDList []uint
	l.svcCtx.DB.Model(group_models.GroupMemberModel{}).
		Where("user_id = ?", req.UserID).
		Select("group_id").Scan(&userGroupIDList)
	//2.子查询查询是否群聊置顶
	column := fmt.Sprintf("(if((select 1 from group_user_top_models where user_id = %d and group_user_top_models.group_id = group_msg_models.group_id), 1,0)) as isTop", req.UserID)

	//3.查哪些聊天记录是被删掉的
	var msgDeleteIDList []uint
	l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).Where("group_id in ?", userGroupIDList).Select("msg_id").Scan(&msgDeleteIDList)

	//4.筛选出我所在的群聊中的聊天记录id，不包括删除的聊天记录id
	query := l.svcCtx.DB.Where("group_id in (?)", userGroupIDList)
	if len(msgDeleteIDList) > 0 {
		query.Where("id not in ?", msgDeleteIDList)
	}
	//5.查询出我的每个群的最新消息，sessionList
	sessionList, count, _ := list_query.ListQuery(l.svcCtx.DB, SessionData{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "isTop desc,newMsgDate desc",
		},

		Table: func() (string, any) {
			return "(?) as u", l.svcCtx.DB.Model(&group_models.GroupMsgModel{}).
				Select("group_id",
					"max(created_at) as newMsgDate",
					column,
					"(select msg_preview from group_msg_models as g where g.group_id = group_id order by g.created_at desc limit 1)as newMsgPreview"). //最新消息预览
				Where(query).
				Group("group_id")
		},
	})

	//6.遍历结果添加到groupIDList中，并查询群信息
	var groupIDList []uint
	for _, data := range sessionList {
		groupIDList = append(groupIDList, data.GroupID)
	}
	var groupListModel []group_models.GroupModel
	l.svcCtx.DB.Find(&groupListModel, groupIDList)
	//7.存储群组信息
	var groupMap = map[uint]group_models.GroupModel{}
	for _, model := range groupListModel {
		groupMap[model.ID] = model
	}

	//8.构造响应
	resp = new(types.GroupSessionListResponse)
	for _, data := range sessionList {
		resp.List = append(resp.List, types.GroupSessionResponse{
			GroupID:       data.GroupID,
			Title:         groupMap[data.GroupID].Title,
			Avatar:        groupMap[data.GroupID].Avatar,
			NewMsgDate:    data.NewMsgDate,
			NewMsgPreview: data.NewMsgPreview,
			IsTop:         data.IsTop,
		})
	}
	resp.Count = int(count)

	fmt.Println(sessionList, count)
	return
}
