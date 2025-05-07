package logic

import (
	"context"
	"errors"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils/set"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupInfoLogic {
	return &GroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupInfoLogic) GroupInfo(req *types.GroupInfoRequest) (resp *types.GroupInfoResponse, err error) {
	//只有群成员才能够查看群详细信息

	//1.根据群id获取群信息，并预加载群成员信息
	var groupModel group_models.GroupModel
	err = l.svcCtx.DB.Preload("MemberList").Take(&groupModel, req.ID).Error
	if err != nil {
		return nil, errors.New("群不存在")
	}

	//2.判断调用者是否是群成员，如果是则获取群成员信息
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}

	//3.构造响应
	resp = &types.GroupInfoResponse{
		GroupID:         groupModel.ID,                                          //群id
		Title:           groupModel.Title,                                       //群名
		Abstract:        groupModel.Abstract,                                    //群头像
		MemberCount:     len(groupModel.MemberList),                             //群成员总数
		Avatar:          groupModel.Avatar,                                      //群简介
		Role:            member.Role,                                            //调用者的角色
		IsProhibition:   groupModel.IsProhibition,                               //是否开启全员禁言
		ProhibitionTime: member.GetProhibitionTime(l.svcCtx.Redis, l.svcCtx.DB), //该用户的群禁言时间
	}
	//4.将群成员id添加到列表中
	//userIDList群主或管理员
	//userAllIDList所有成员
	var userIDList []uint32
	var userAllIDList []uint32
	for _, model := range groupModel.MemberList {
		if model.Role == 1 || model.Role == 2 {
			userIDList = append(userIDList, uint32(model.UserID))
		}
		userAllIDList = append(userAllIDList, uint32(model.UserID))
	}
	//5.获取群主、管理员列表的用户信息
	userListResponse, err := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	if err != nil {
		return
	}
	var creator types.UserInfo
	var adminList = make([]types.UserInfo, 0)

	//6.算在线用户总数
	userOnlineResponse, err := l.svcCtx.UserRpc.UserOlineList(context.Background(), &user_rpc.UserOlineListRequest{})
	if err == nil {
		// 算群成员和总的在线人数成员,取交集
		slice := set.Intersect(userOnlineResponse.UserIdList, userAllIDList)
		resp.MemberOnlineCount = len(slice)
	}

	//7.遍历群成员列表，判断群成员角色，如果是群主或管理，则记录信息
	for _, model := range groupModel.MemberList {
		if model.Role == 3 {
			continue
		}
		userInfo := types.UserInfo{
			UserID:   model.UserID,
			Avatart:  userListResponse.UserInfo[uint32(model.UserID)].Avatar,
			Nickname: userListResponse.UserInfo[uint32(model.UserID)].NickName,
		}
		if model.Role == 1 {
			creator = userInfo
			continue
		}
		if model.Role == 2 {
			adminList = append(adminList, userInfo)
		}
	}
	resp.Creator = creator
	resp.AdminList = adminList

	return
}
