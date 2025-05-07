package logic

import (
	"context"
	"errors"
	"fmt"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils/set"
	"strings"
	"time"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupCreateLogic {
	return &GroupCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupCreateLogic) GroupCreate(req *types.GroupCreateRequest) (resp *types.GroupCreateResponse, err error) {
	//1.初始化群聊结构
	var groupModel = group_models.GroupModel{
		Creator:      req.UserID, //自己创建的群,自己就是群主,
		IsSearch:     false,
		Abstract:     fmt.Sprintf("本群创建于%s:群主很懒,什么都没有留下", time.Now().Format("2006-01-02")),
		Verification: 2,
		Size:         50,
	}
	//groupUserList群成员列表，一定包含自己
	var groupUserList = []uint{req.UserID}
	//2.根据创建模式创建群聊
	switch req.Mode {
	case 1: // 直接创建模式
		if req.Name == "" {
			return nil, errors.New("群名不能为空")
		}
		if req.Size >= 1000 {
			return nil, errors.New("群规模错误")
		}
		groupModel.Title = req.Name
		groupModel.Size = req.Size
		groupModel.IsSearch = req.IsSearch
	case 2: // 选人创建模式
		if len(req.UserIDList) == 0 {
			return nil, errors.New("没有要选择的好友")
		}
		// 去算选择的用户昵称,是不是超过最大长度
		// 群名是32
		//调用户信息列表
		var userIDList = []uint32{uint32(req.UserID)} // 先把自己放进去
		for _, u := range req.UserIDList {
			userIDList = append(userIDList, uint32(u))
			groupUserList = append(groupUserList, u)
		}

		//查询好友列表信息
		userFriendResponse, err := l.svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
			User: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		//提取好友列表id
		var friendIDList []uint
		for _, i2 := range userFriendResponse.FriendList {
			friendIDList = append(friendIDList, uint(i2.UserId))
		}
		//判断它们两个是不是一致的，都是好友才可以选人加群
		slice := set.Difference(req.UserIDList, friendIDList)
		if len(slice) != 0 {
			return nil, errors.New("选择的好友列表中有人不是你的好友")
		}

		//查询用户列表信息
		userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
			UserIdList: userIDList,
		})
		if err1 != nil {
			logx.Error(err1)
			return nil, errors.New("用户服务错误")
		}
		//初始化群标题，根据名字创建
		// 去算这个昵称的长度 算到第几个人的时候会大于32
		var nameList []string
		for _, info := range userListResponse.UserInfo {
			if len(strings.Join(nameList, "、")) >= 32 {
				break
			}
			nameList = append(nameList, info.NickName)
		}
		groupModel.Title = strings.Join(nameList, "、") + "群聊"
	default:
		return nil, errors.New("不支持模式")
	}
	// 群头像
	// 1.默认头像 2.文字头像
	groupModel.Avatar = string([]rune(groupModel.Title)[0])
	//创建群
	err = l.svcCtx.DB.Create(&groupModel).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("创建群组失败")
	}

	//添加群成员信息
	var members []group_models.GroupMemberModel
	for i, u := range groupUserList {
		memberModel := group_models.GroupMemberModel{
			GroupID: groupModel.ID,
			UserID:  u,
			Role:    3,
		}
		if i == 0 {
			memberModel.Role = 1
		}
		members = append(members, memberModel)
	}

	err = l.svcCtx.DB.Create(&members).Error
	if err != nil {
		logx.Error(err)
		return nil, errors.New("群成员添加失败")
	}
	return
}
