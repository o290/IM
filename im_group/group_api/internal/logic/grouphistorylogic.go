package logic

import (
	"context"
	"errors"
	"server/common/list_query"
	"server/common/models"
	"server/common/models/ctype"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils"
	"time"

	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGroupHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupHistoryLogic {
	return &GroupHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

type HistoryResponse struct {
	UserID         uint          `json:"userID"`
	UserNickname   string        `json:"userNickname"`
	UserAvatar     string        `json:"userAvatar"`
	Msg            ctype.Msg     `json:"msg"`
	ID             uint          `json:"ID"`
	MsgType        ctype.MsgType `json:"msgType"`
	CreatedAt      time.Time     `json:"createdAt"`
	IsMe           bool          `json:"isMe"`
	MemberNickname string        `json:"memberNickname"` //群用户备注
}
type HistoryListResponse struct {
	List  []HistoryResponse `json:"list"`
	Count int               `json:"count"`
}

func (l *GroupHistoryLogic) GroupHistory(req *types.GroupHistoryRequest) (resp *HistoryListResponse, err error) {
	//1.判断调用者是否是群成员
	var member group_models.GroupMemberModel
	err = l.svcCtx.DB.Take(&member, "group_id = ? and user_id = ?", req.ID, req.UserID).Error
	if err != nil {
		return nil, errors.New("该用户不是群成员")
	}
	//2.查询调用者删除了哪些聊天记录
	var msgIDList []uint
	l.svcCtx.DB.Model(group_models.GroupUserMsgDeleteModel{}).
		Where("group_id = ? and user_id = ?", req.ID, req.UserID).
		Select("msg_id").Scan(&msgIDList)

	//3.过滤掉删除的记录
	var query = l.svcCtx.DB.Where("")
	if len(msgIDList) > 0 {
		query.Where("id not in ?", msgIDList)
	}
	//4.查询群聊id的聊天记录
	groupsMsgList, count, _ := list_query.ListQuery(l.svcCtx.DB, group_models.GroupMsgModel{GroupID: req.ID}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc",
		},
		Where:   query,
		Preload: []string{"GroupMemberModel"},
	})

	//5.遍历群聊消息，将发送方的id保存，然后去重，查询发送方id的用户信息
	var userIDList []uint32
	for _, model := range groupsMsgList {
		userIDList = append(userIDList, uint32(model.SendUserID))
	}
	userIDList = utils.DeduplicationList(userIDList)
	userListResponse, err1 := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})

	//6.构造响应
	var list = make([]HistoryResponse, 0)
	for _, model := range groupsMsgList {
		info := HistoryResponse{
			UserID:    model.SendUserID,
			Msg:       model.Msg,
			ID:        model.ID,
			MsgType:   model.MsgType,
			CreatedAt: model.CreatedAt,
		}
		if model.GroupMemberModel != nil {
			info.MemberNickname = model.GroupMemberModel.MemberNickname
		}
		if err1 == nil {
			info.UserNickname = userListResponse.UserInfo[uint32(info.UserID)].NickName
			info.UserAvatar = userListResponse.UserInfo[uint32(info.UserID)].Avatar
		}
		//判断是不是自己发的
		if req.UserID == info.UserID {
			info.IsMe = true
		}
		list = append(list, info)
	}
	resp = new(HistoryListResponse)
	resp.List = list
	resp.Count = int(count)

	return
}
