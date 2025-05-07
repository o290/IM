package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"server/common/list_query"
	"server/common/models"
	"server/common/models/ctype"
	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"
	"server/im_user/user_rpc/types/user_rpc"
)

type GroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}
type Data struct {
	GroupID        uint   `gorm:"column:group_id"`
	UserID         uint   `gorm:"column:user_id"`
	Role           int8   `gorm:"column:role"`
	CreatedAt      string `gorm:"column:created_at"`
	MemberNickname string `gorm:"column:member_nick_name"`
	NewMsgDate     string `gorm:"column:new_msg_date"`
}

func NewGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupMemberLogic {
	return &GroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupMemberLogic) GroupMember(req *types.GroupMemberRequest) (resp *types.GroupMemberResponse, err error) {
	//1.检查是否支持排序模式
	switch req.Sort {
	case "new_msg_date desc", "new_msg_date asc": //按照最新发言排序
	case "role asc": // 按照角色升序
	case "created_at desc", "created_at asc": //按照进群时间排
	default:
		return nil, errors.New("不支持的排序模式")
	}
	//2.查询用户在该群聊中个的最新发言时间
	// 构建子查询的 SQL 表达式，使用 MAX 函数
	subQuerySQL := "(SELECT MAX(group_msg_models.created_at) FROM group_msg_models WHERE group_msg_models.group_id =? AND group_msg_models.send_user_id = user_id) AS new_msg_date"
	//subQuerySQL := "(SELECT group_msg_models.created_at FROM group_msg_models WHERE group_msg_models.group_id =? AND group_msg_models.send_user_id = user_id) AS new_msg_date"

	//3.构建主查询，查询用户在这个群聊里的
	//select
	//	group_id,user_id,role,created_at,
	//		(SELECT MAX(g.created_at)
	//	FROM group_msg_models g
	//	WHERE
	//	g.group_id =6 AND g.send_user_id = 1)as new_msg
	//	from group_member_models
	//	where group_id=6
	query := l.svcCtx.DB.Table("group_member_models").
		Where("group_id =?", req.ID).
		Select("group_id, user_id, role, created_at, member_nickname, "+subQuerySQL, req.ID)

	//4.调用 ListQuery 进行查询
	//select u.group_id,u.user_id,u.role,u.created_at,u.new_msg,u.member_nickname
	//	from
	//	(select
	//		group_id,user_id,role,created_at,member_nickname,
	//			(SELECT MAX(g.created_at)
	//		FROM group_msg_models g
	//		WHERE
	//		g.group_id =6 AND g.send_user_id = 1)as new_msg
	//		from group_member_models
	//		where group_id=6)as u
	//		order by created_at desc
	//		limit 10 offset 0;
	//memberList是每个用户在该群聊中最新的消息记录
	memberList, count, err := list_query.ListQuery(l.svcCtx.DB, Data{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  req.Sort,
		},
		Table: func() (string, any) {
			return "(?) as u", query
		},
	})
	//5.从memberList获取用户id，然后根据用户id列表获取用户信息
	var userIDList []uint32
	for _, data := range memberList {
		userIDList = append(userIDList, uint32(data.UserID))
	}

	var userInfoMap = map[uint]ctype.UserInfo{}
	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	if err == nil {
		for u, info := range userListResponse.UserInfo {
			userInfoMap[uint(u)] = ctype.UserInfo{
				ID:       uint(u),
				NickName: info.NickName,
				Avatar:   info.Avatar,
			}
		}
	} else {
		logx.Error(err)
	}

	//6.获取用户的在线状态
	var userOnlineMap = map[uint]bool{}
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOlineList(context.Background(), &user_rpc.UserOlineListRequest{})
	if err == nil {
		for _, u := range userOnlineResponse.UserIdList {
			userOnlineMap[uint(u)] = true
		}
	} else {
		logx.Error(err)
	}

	//7.返回响应
	resp = new(types.GroupMemberResponse)
	for _, data := range memberList {
		resp.List = append(resp.List, types.GroupMemberInfo{
			UserID:         data.UserID,
			UserNickname:   userInfoMap[data.UserID].NickName,
			Avatar:         userInfoMap[data.UserID].Avatar,
			IsOnline:       userOnlineMap[data.UserID],
			Role:           data.Role,
			MemberNickname: data.MemberNickname,
			CreatedAt:      data.CreatedAt,
			NewMsgDate:     data.NewMsgDate,
		})
	}
	resp.Count = int(count)
	fmt.Println(count)
	return
}
