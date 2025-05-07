package logic

import (
	"context"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupFriendsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupFriendsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupFriendsListLogic {
	return &GroupFriendsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupFriendsListLogic) GroupFriendsList(req *types.GroupFriendsListRequest) (resp *types.GroupFriendsListResponse, err error) {
	// 我的好友哪些在这个群里面

	//1.查询调用者的好友列表
	friendResponse, err := l.svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
		User: uint32(req.UserID),
	})
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	//2.查询这个群的群成员列表
	var memberList []group_models.GroupMemberModel
	l.svcCtx.DB.Find(&memberList, "group_id = ?", req.ID)

	//3.标识哪些用户在群中
	var memberMap = map[uint]bool{}
	for _, model := range memberList {
		memberMap[model.UserID] = true
	}
	//4.构造响应，遍历用户的好友列表
	resp = new(types.GroupFriendsListResponse)
	count := 0
	//fmt.Println(friendResponse.FriendList)
	for _, info := range friendResponse.FriendList {
		//fmt.Println(info)
		resp.List = append(resp.List, types.GroupFriendsResponse{
			UserId:    uint(info.UserId),
			Avatar:    info.Avatar,
			Nickname:  info.NickName,
			IsInGroup: memberMap[uint(info.UserId)],
		})
		if _, ok := memberMap[uint(info.UserId)]; ok {
			count++
		}
	}
	resp.Count = count
	return
}
