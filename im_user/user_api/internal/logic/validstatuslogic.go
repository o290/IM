package logic

import (
	"context"
	"encoding/json"
	"errors"
	"server/common/models/ctype"
	"server/im_chat/chat_rpc/chat"
	"server/im_user/user_models"

	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidStatusLogic {
	return &ValidStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidStatusLogic) ValidStatus(req *types.FriendValidStatusRequest) (resp *types.FriendValidStatusResponse, err error) {
	//1.查询验证记录
	var friendVerify user_models.FriendVerifyModel
	err = l.svcCtx.DB.Take(&friendVerify, "id=? and rev_user_id=?", req.VerifyID, req.UserID).Error
	//自己是接收方
	if err != nil {
		return nil, errors.New("记录不存在")
	}
	//2.判断是否是未操作状态
	if friendVerify.RevStatus != 0 {
		return nil, errors.New("不可更改状态")
	}
	//3.根据请求体中的状态分别操作
	switch req.Status {
	case 1: //同意
		friendVerify.RevStatus = 1
		//往好友表里添加
		l.svcCtx.DB.Create(&user_models.FriendModel{
			SendUserID: friendVerify.SendUserID,
			RevUserID:  friendVerify.RevUserID,
		})

		msg := ctype.Msg{
			Type: ctype.TextMsgType,
			TextMsg: &ctype.TextMsg{
				Content: "我们已经是好友了，开始聊天吧",
			},
		}
		byteData, _ := json.Marshal(msg)
		//给对方发送消息
		_, err = l.svcCtx.ChatRpc.UserChat(context.Background(), &chat.UserChatRequest{
			SendUserId: uint32(friendVerify.SendUserID),
			RevUserId:  uint32(friendVerify.RevUserID),
			Msg:        byteData,
			SystemMsg:  nil,
		})
		if err != nil {
			logx.Error(err)

		}
	case 2: //拒绝
		friendVerify.RevStatus = 2
	case 3: //忽略
		friendVerify.RevStatus = 3
	case 4: //删除
		//一条验证记录是两个人看的
		l.svcCtx.DB.Delete(&friendVerify)
		return nil, nil
	}
	l.svcCtx.DB.Save(&friendVerify)
	return
}
