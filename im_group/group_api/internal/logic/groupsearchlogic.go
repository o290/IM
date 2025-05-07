package logic

import (
	"context"
	"server/common/list_query"
	"server/common/models"
	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils/set"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupSearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupSearchLogic {
	return &GroupSearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupSearchLogic) GroupSearch(req *types.GroupSearchRequest) (resp *types.GroupSearchListResponse, err error) {
	//1.根据群id或群名查询群信息
	//isSearch为false为不能被搜索
	groups, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupModel{}, list_query.Option{
		//分页查询
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		//预加载群成员信息
		Preload: []string{"MemberList"},
		//群可以被搜索到，且群id相等或者群名模糊匹配
		Where: l.svcCtx.DB.Debug().Where("is_search=1 and(id = ? or title LIKE ?)", req.Key, "%"+req.Key+"%"),
	})

	//2.获取在线用户列表
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOlineList(context.Background(), &user_rpc.UserOlineListRequest{})
	var userOnlineIDList []uint
	if err == nil {
		//称之为:服务降级,如果用户rpc方法挂了,只是页面上看到在线人数是0而已,不会影响这个群搜索功能
		for _, u := range userOnlineResponse.UserIdList {
			userOnlineIDList = append(userOnlineIDList, uint(u))
		}
	}

	//3.组装搜索结果
	resp = new(types.GroupSearchListResponse)
	//遍历查询结果groups
	for _, group := range groups {
		//groupMemberIdList存储当前群组的所有成员的用户 ID
		var groupMemberIdList []uint
		//标当req.UserID是否在当前群组中
		var isInGroup bool
		//遍历当前群组的成员列表
		for _, model := range group.MemberList {
			groupMemberIdList = append(groupMemberIdList, model.UserID)
			if model.UserID == req.UserID {
				isInGroup = true
			}
		}
		//封装群组信息
		resp.List = append(resp.List, types.GroupSearchResponse{
			GroupID:         group.ID,
			Title:           group.Title,
			Abstract:        group.Abstract,
			Avatar:          group.Avatar,
			UserCount:       len(group.MemberList),
			UserOnlineCount: len(set.Intersect(groupMemberIdList, userOnlineIDList)), //这个群的在线用户总数
			IsInGroup:       isInGroup,
		})
	}
	resp.Count = int(count)

	return
}
