package logic

import (
	"context"
	"fmt"
	"server/im_chat/chat_models"

	"server/im_chat/chat_api/internal/svc"
	"server/im_chat/chat_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChatDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChatDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatDeleteLogic {
	return &ChatDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChatDeleteLogic) ChatDelete(req *types.ChatDeleteRequest) (resp *types.ChatDeleteResponse, err error) {
	//1.从chat_model表中查询请求id中的聊天记录数据绑定在chatList中
	//IDList是包含聊天记录id的集合
	var chatList []chat_models.ChatModel
	l.svcCtx.DB.Find(&chatList, req.IDList)

	//2.从chat_delete表中查询请求id中的聊天记录数据绑定在useDeleteChatList中，并存储在map中
	var useDeleteChatList []chat_models.UserChatDeleteModel
	l.svcCtx.DB.Find(&useDeleteChatList, req.IDList)
	// struct{} 不占用内存空间
	chatDeleteMap := map[uint]struct{}{}
	for _, model := range useDeleteChatList {
		chatDeleteMap[model.ChatID] = struct{}{}
	}

	//3.遍历chatList，判断是否存在于chatDeleteMap中，不存在则添加至deleteChatIDList终
	var deleteChatIDList []chat_models.UserChatDeleteModel
	if len(chatList) > 0 {
		for _, model := range chatList {
			//不是自己的聊天记录
			if !(model.SendUserID == req.UserID || model.RevUserID == req.UserID) {
				fmt.Println("不是自己的聊天记录", model.ID)
				continue
			}
			//已经删过的聊天记录
			_, ok := chatDeleteMap[model.ID]
			if ok {
				fmt.Println("已经删过了", model.ID)
				continue
			}
			fmt.Println(req.UserID, model.ID)
			deleteChatIDList = append(deleteChatIDList, chat_models.UserChatDeleteModel{
				UserID: req.UserID,
				ChatID: model.ID,
			})
		}
	}
	//4.添加至chat_delete中
	if len(deleteChatIDList) > 0 {
		l.svcCtx.DB.Create(&deleteChatIDList)
	}

	logx.Infof("已删除聊天记录 %d 条", len(deleteChatIDList))
	return
}
