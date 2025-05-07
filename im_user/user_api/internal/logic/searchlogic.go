package logic

import (
	"context"
	"fmt"
	"server/common/list_query"
	"server/common/models"
	"server/im_user/user_models"

	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchRequest) (resp *types.SearchResponse, err error) {
	//1.搜索所有符合条件的用户
	users, count, _ := list_query.ListQuery(l.svcCtx.DB, user_models.UserConfModel{
		//搜索在线用户
		Online: req.Online,
	}, list_query.Option{
		//分页查询
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
		},
		//预加载
		Preload: []string{"UserModel"},
		//获取用户基本信息
		Joins: "left join user_models um on um.id = user_conf_models.user_id",
		Where: l.svcCtx.DB.Where("(user_conf_models.search_user<>0 or user_conf_models.search_user is not null) and (user_conf_models.search_user = 1 and um.id = ?) or (user_conf_models.search_user=2 and (um.id=? or um.nickname like ?))", req.Key, req.Key, fmt.Sprintf("%%%s%%", req.Key)),
	})
	//先筛选出有效的search_user字段，然后再根据search_user进一步筛选
	//(user_conf_models.search_user<>0 or user_conf_models.search_user is not null)筛选出search_user非空且不为0，即查找用户设置了允许比如查询的用户
	//and
	//(user_conf_models.search_user = 1 and um.id = ?)如果search_user为1，要求传入的id与req.Key相同
	//or
	//(user_conf_models.search_user=2 and (um.id=? or um.nickname like ?))如果search_user为2，要求传入的id与req.Key相同或者与nickname相同
	//2.查当前用户的好友列表
	var friend user_models.FriendModel
	friends := friend.Friends(l.svcCtx.DB, req.UserID)
	userMap := map[uint]bool{}
	//userMap存储是是用户的好友id
	for _, model := range friends {
		if model.SendUserID == req.UserID {
			userMap[model.RevUserID] = true
		} else {
			userMap[model.SendUserID] = true
		}
	}
	//3.查询所有用户
	list := make([]types.SearchInfo, 0)
	for _, user := range users {
		list = append(list, types.SearchInfo{
			UserID:   user.UserID,
			NickName: user.UserModel.Nickname,
			Abstract: user.UserModel.Abstract,
			Avatar:   user.UserModel.Avatar,
			//如果是好友则标记
			IsFriend: userMap[user.UserID],
		})
	}

	return &types.SearchResponse{List: list, Count: count}, nil
}
