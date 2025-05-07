package logic

import (
	"context"
	"server/common/list_query"
	"server/common/models"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupValidListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupValidListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupValidListLogic {
	return &GroupValidListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupValidListLogic) GroupValidList(req *types.GroupValidListRequest) (resp *types.GroupValidListResponse, err error) {
	var groupIDList []uint //我管理的群
	//1.查询在我加入的群我是管理员或群主的群id集合
	l.svcCtx.DB.Model(group_models.GroupMemberModel{}).Where("user_id = ? and (role = 1 or role = 2)", req.UserID).Select("group_id").Scan(&groupIDList)
	//2.查询在群id集合中的验证列表(别人发的)或者是我发起的验证记录
	groups, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupVerifyModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Preload: []string{"GroupModel"},
		Where:   l.svcCtx.DB.Where("group_id in ? or user_id = ?", groupIDList, req.UserID),
	})

	//3.将查询结果中出现的用户加入到userIDList中
	var userIDList []uint32
	for _, group := range groups {
		userIDList = append(userIDList, uint32(group.UserID))
	}

	//4.获取用户列表信息
	userList, err1 := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	//5.构造响应
	resp = new(types.GroupValidListResponse)
	resp.Count = int(count)
	for _, group := range groups {
		info := types.GroupValidInfoResponse{
			ID:                 group.ID,
			GroupID:            group.GroupID,
			UserID:             group.UserID,
			Status:             group.Status,
			AdditionalMessages: group.AdditionalMessage,
			Title:              group.GroupModel.Title,
			CreatedAt:          group.CreatedAt.String(),
			Type:               group.Type,
		}
		if group.VerificationQuestion != nil {
			info.VerificationQuestion = &types.VerificationQuestion{
				Problem1: group.VerificationQuestion.Problem1,
				Problem2: group.VerificationQuestion.Problem2,
				Problem3: group.VerificationQuestion.Problem3,
				Answer1:  group.VerificationQuestion.Answer1,
				Answer2:  group.VerificationQuestion.Answer2,
				Answer3:  group.VerificationQuestion.Answer3,
			}
		}
		if err1 == nil {
			info.UserNickname = userList.UserInfo[uint32(info.UserID)].NickName
			info.UserAvatar = userList.UserInfo[uint32(info.UserID)].Avatar
		}
		resp.List = append(resp.List, info)
	}
	return
}
