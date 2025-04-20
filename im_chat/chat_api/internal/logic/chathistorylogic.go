package logic

import (
	"context"
	"errors"
	"server/common/list_query"
	"server/common/models"
	"server/common/models/ctype"
	"server/im_chat/chat_api/internal/svc"
	"server/im_chat/chat_api/internal/types"
	"server/im_chat/chat_models"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatHistoryLogic {
	return &ChatHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 自己构造响应
type ChatHistory struct {
	ID        uint             `json:"id"`
	SendUser  ctype.UserInfo   `json:"sendUser"`
	RevUser   ctype.UserInfo   `json:"revUser"`
	IsMe      bool             `json:"isMe"`       //哪条消息是我发的
	CreatedAt string           `json:"created_at"` //消息时间
	Msg       ctype.Msg        `json:"msg"`        //消息
	SystemMsg *ctype.SystemMsg `json:"systemMsg"`  //系统消息
}
type ChatHistoryResponse struct {
	List  []ChatHistory `json:"list"`
	Count int64         `json:"count"`
}

func (l *ChatHistoryLogic) ChatHistory(req *types.ChatHistoryRequest) (resp *ChatHistoryResponse, err error) {
	//1.判断是否是好友
	if req.UserID != req.FriendID {
		res, err := l.svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
			User2: uint32(req.UserID),
			User1: uint32(req.FriendID),
		})
		if err != nil {
			return nil, err
		}
		if !res.IsFriend {
			return nil, errors.New("你们还不是好友")
		}
	}

	//2.查询自己和好友的聊天记录
	//select * from chat_models
	//	where
	//	send_user_id=1 and rev_user_id=2 or send_user_id=2 and rev_user_id=1
	//	and id not in (
	//	select chat_id from user_chat_delete_models
	//		where user_id=1
	//		)
	//		order by created_at desc
	//		limit 1 offset 1;
	chatList, count, _ := list_query.ListQuery(l.svcCtx.DB, chat_models.ChatModel{}, list_query.Option{
		//分页查询
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.Limit,
			Sort:  "created_at desc", //根据创建时间降序
		},
		Where: l.svcCtx.DB.Where("((send_user_id = ? and rev_user_id = ?)or(send_user_id = ? and rev_user_id = ?) )and id not in (select chat_id from user_chat_delete_models where user_id=?)",
			req.UserID, req.FriendID, req.FriendID, req.UserID, req.UserID),
	})

	//3.遍历聊天记录列表，收集所有的用户id，对userIDList去重
	var userIDList []uint32
	for _, model := range chatList {
		userIDList = append(userIDList, uint32(model.SendUserID))
		userIDList = append(userIDList, uint32(model.RevUserID))
	}
	//去重
	userIDList = utils.DeduplicationList(userIDList)

	//4.调用用户服务的rpc方法，获取用户列表信息{用户 id：{用户信息}}
	response, err := l.svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
		UserIdList: userIDList,
	})
	if err != nil {
		logx.Error(err)
		return nil, errors.New("用户服务错误")
	}

	//5.组装聊天记录响应
	var list = make([]ChatHistory, 0)
	for _, model := range chatList {
		sendUser := ctype.UserInfo{
			ID:       model.SendUserID,
			NickName: response.UserInfo[uint32(model.SendUserID)].NickName,
			Avatar:   response.UserInfo[uint32(model.SendUserID)].Avatar,
		}
		revUser := ctype.UserInfo{
			ID:       model.RevUserID,
			NickName: response.UserInfo[uint32(model.RevUserID)].NickName,
			Avatar:   response.UserInfo[uint32(model.RevUserID)].Avatar,
		}
		info := ChatHistory{
			ID:        model.ID,
			CreatedAt: model.CreatedAt.String(),
			SendUser:  sendUser,
			RevUser:   revUser,
			Msg:       model.Msg,
			SystemMsg: model.SystemMsg,
		}
		//IsMe是不是自己发的
		if info.SendUser.ID == req.UserID {
			info.IsMe = true
		}
		list = append(list, info)
	}
	resp = &ChatHistoryResponse{
		List:  list,
		Count: count,
	}
	return
}
