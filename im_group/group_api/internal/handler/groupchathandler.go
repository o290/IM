package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"gorm.io/gorm"
	"net/http"
	"server/common/models/ctype"
	"server/common/response"
	"server/common/service/redis_service"
	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"
	"server/im_group/group_models"
	"server/im_user/user_rpc/types/user_rpc"
	"time"
)

type UserWsInfo struct {
	UserInfo    ctype.UserInfo             //用户信息
	WsClientMap map[string]*websocket.Conn //这个用户管理的所有客户端
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

type ChatRequest struct {
	GroupID uint      `json:"groupID"`
	Msg     ctype.Msg `json:"msg"`
}
type ChatResponse struct {
	UserID         uint          `json:"userID"`
	UserNickname   string        `json:"userNickname"`
	UserAvatar     string        `json:"userAvatar"`
	IsMe           bool          `json:"isMe"`
	Msg            ctype.Msg     `json:"msg"`
	ID             uint          `json:"id"`
	MsgTyp         ctype.MsgType `json:"msgTyp"`
	CreatedAt      time.Time     `json:"created_at"`
	MemberNickname string        `json:"memberNickname"` //群用户备注
}

func groupChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//一.初始化与处理请求
		//1.从header中解析请求，得到请求体req
		var req types.GroupChatRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		//解析handler
		if err := httpx.ParseHeaders(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		//2.创建一个 websocket.Upgrader 实例用于将 HTTP 连接升级为 WebSocket 连接
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}
		//3.将http升级为websocket
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		//二.用户连接管理
		//1.RemoteAddr获取当前连接的远程地址，并转换字符串
		addr := conn.RemoteAddr().String()
		logx.Infof("用户建立ws连接 %s", addr)
		//2.处理websocket连接关闭
		defer func() {
			conn.Close()
			//获取当前用户的在线信息，在线则退出连接
			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				//删除退出的ws信息
				delete(userWsInfo.WsClientMap, addr)
			}
			//检查该用户的所有连接是否全部关闭
			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				//全退完了，删除对应缓存
				delete(UserOnlineWsMap, req.UserID)
			}
		}()
		//获取用户基本信息
		baseInfoResponse, err := svcCtx.UserRpc.UserBaseInfo(context.Background(), &user_rpc.UserBaseInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return

		}

		//将服务返回的 JSON 数据解析到 userInfo 结构体中
		userInfo := ctype.UserInfo{
			ID:       req.UserID,
			NickName: baseInfoResponse.NickName,
			Avatar:   baseInfoResponse.Avatar,
		}
		//4.更新用户的在线信息
		//查找用户是否在线，若不在线，表示这是用户的首次连接
		userWsInfo, ok := UserOnlineWsMap[req.UserID]
		if !ok {
			//用户首次连接，存入map中
			userWsInfo = &UserWsInfo{
				UserInfo: userInfo,
				WsClientMap: map[string]*websocket.Conn{
					addr: conn,
				},
			}
			UserOnlineWsMap[req.UserID] = userWsInfo
		}
		//检查当前客户端地址是否已经存在于用户的 WsClientMap 中
		_, ok1 := userWsInfo.WsClientMap[addr]
		//如果不存在，说明用户是第二次或更多次连接
		if !ok1 {
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
		}

		//四.循环从websocket连接conn中读取消息，消息预处理，根据不同类型消息处理，然后存储并发送消息
		for {
			//读取消息：消息类型 消息 错误
			_, p, err1 := conn.ReadMessage()
			if err1 != nil {
				//用户断开聊天
				break
			}
			//解析消息
			var request ChatRequest
			err = json.Unmarshal(p, &request)
			if err != nil {
				SendTipErrMsg(conn, "参数解析失败")
				continue
			}

			// 判断自己是不是这个群的成员
			var member group_models.GroupMemberModel
			err = svcCtx.DB.Preload("GroupModel").Take(&member, "group_id = ? and user_id = ?", request.GroupID, req.UserID).Error
			if err != nil {
				// 自己不是群的成员
				SendTipErrMsg(conn, "你还不是群成员")
				continue
			}

			if member.GroupModel.IsProhibition {
				// 开启了全员禁言
				SendTipErrMsg(conn, "当前群正在全员禁言中")
				continue
			}
			// 我是不是被禁言了
			if member.GetProhibitionTime(svcCtx.Redis, svcCtx.DB) != nil {
				SendTipErrMsg(conn, "当前用户正在禁言中")
				continue
			}

			//根据消息类型进行处理
			switch request.Msg.Type {
			case ctype.WithdrawMsgType:
				//校验
				withdrawMsg := request.Msg.WithdrawMsg
				if request.Msg.WithdrawMsg == nil {
					SendTipErrMsg(conn, "撤回消息格式错误")
					continue
				}
				if withdrawMsg.MsgID == 0 {
					SendTipErrMsg(conn, "撤回消息id为空")
					continue
				}
				// 去找消息
				var groupMsg group_models.GroupMsgModel
				err = svcCtx.DB.Take(&groupMsg, "group_id=? and id = ?", request.GroupID, withdrawMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "原消息不存在")
					continue
				}

				// 原消息不能是撤回消息
				if groupMsg.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue

				}
				// 要去拿我在这个群的角色
				if member.Role == 3 {
					// 如果是自己撤自己的 并且自己是普通用户
					if req.UserID != groupMsg.SendUserID {
						SendTipErrMsg(conn, "普通用户只能撤回自己的消息")
						continue
					}
					// 要判断时间是不是大于了2分钟
					now := time.Now()
					if now.Sub(groupMsg.CreatedAt) > 2*time.Minute {
						SendTipErrMsg(conn, "只能撤回两分钟以内的消息")
						continue
					}
				}
				// 查这个消息的用户,在这个群里的角色
				var msgUserRole int8 = 3
				err = svcCtx.DB.Model(group_models.GroupMemberModel{}).
					Where("group_id = ? and user_id = ?", request.GroupID, groupMsg.SendUserID).
					Select("role").
					Scan(&msgUserRole).Error
				//这里有可能查不到 原因是这个消息的用户退群了,那么也是可以撤回的

				// 如果是管理员撤回 它能撤自己和用户的,没有时间限制
				if member.Role == 2 {
					// 不能撤群主和别的管理员
					if msgUserRole == 1 || (msgUserRole == 2 && groupMsg.SendUserID != req.UserID) {
						SendTipErrMsg(conn, "管理员只能撤回自己或者普通用户的消息")
						continue
					}
				}

				// 如果是群主,那就能撤管理员和用户的

				//代表消息可以茶会
				//修改原消息
				var content = "撤回了一条消息"
				content = "你" + content
				// 前端可以判断,这个消息如果不是isMe,就可以把你替换成对方的昵称
				originMsg := groupMsg.Msg
				// 这里可能会出现循环引用,所以拷贝了这个值,并且把撤回消息置空了
				originMsg.WithdrawMsg = nil
				svcCtx.DB.Model(&groupMsg).Updates(group_models.GroupMsgModel{
					MsgPreview: "[撤回消息]-" + content,
					MsgType:    ctype.WithdrawMsgType,
					Msg: ctype.Msg{
						Type: ctype.WithdrawMsgType,
						WithdrawMsg: &ctype.WithdrawMsg{
							Content:   content,
							MsgID:     request.Msg.WithdrawMsg.MsgID,
							OriginMsg: &originMsg,
						},
					},
				})
			case ctype.ReplyMsgType:
				//回复消息
				//先校验
				if request.Msg.ReplyMsg == nil || request.Msg.ReplyMsg.MsgID == 0 {
					SendTipErrMsg(conn, "回复消息id必填")
					return
				}
				// 找这个原消息
				var msgModel group_models.GroupMsgModel
				err = svcCtx.DB.Take(&msgModel, "group_id=? and id=?", request.GroupID, request.Msg.ReplyMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				//不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}

				userBaseInfo, err4 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)
				if err4 != nil {
					logx.Error(err)
					SendTipErrMsg(conn, err4.Error())
					continue
				}
				request.Msg.ReplyMsg.Msg = &msgModel.Msg
				request.Msg.ReplyMsg.UserID = msgModel.SendUserID
				request.Msg.ReplyMsg.UserNickName = userBaseInfo.NickName
				request.Msg.ReplyMsg.OriginMsgDate = msgModel.CreatedAt
			case ctype.QuoteMsgType:
				//先校验
				if request.Msg.QuoteMsg == nil || request.Msg.QuoteMsg.MsgID == 0 {
					SendTipErrMsg(conn, "引用消息id必填")
					return
				}
				// 找这个原消息
				var msgModel group_models.GroupMsgModel
				err = svcCtx.DB.Take(&msgModel, "group_id=? and id=?", request.GroupID, request.Msg.QuoteMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				//不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}

				userBaseInfo, err4 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)
				if err4 != nil {
					logx.Error(err)
					SendTipErrMsg(conn, err4.Error())
					continue
				}
				request.Msg.QuoteMsg.Msg = &msgModel.Msg
				request.Msg.QuoteMsg.UserID = msgModel.SendUserID
				request.Msg.QuoteMsg.UserNickName = userBaseInfo.NickName
				request.Msg.QuoteMsg.OriginMsgDate = msgModel.CreatedAt
			}
			msgID := insertMsg(svcCtx.DB, conn, member, request.Msg)
			// 遍历这个用户列表,去找ws的客户端，把消息发送给所有人
			sendGroupOlineUserMSg(
				svcCtx.DB,
				member,
				request.Msg,
				msgID)
		}
	}
}

func insertMsg(db *gorm.DB, conn *websocket.Conn, member group_models.GroupMemberModel, msg ctype.Msg) uint {
	switch msg.Type {
	case ctype.WithdrawMsgType:
		fmt.Println("撒回消息自己是不入库的")
		return 0
	}
	groupMsg := group_models.GroupMsgModel{
		GroupID:       member.GroupID,
		SendUserID:    member.UserID,
		GroupMemberID: member.ID,
		MsgType:       msg.Type,
		Msg:           msg,
	}
	groupMsg.MsgPreview = groupMsg.MsgPreviewMethod()
	err := db.Create(&groupMsg).Error
	if err != nil {
		logx.Error(err)
		SendTipErrMsg(conn, "消息保存失败")
		return 0
	}
	return groupMsg.ID
}

// sendGroupOlineUserMSg 给这个群的用户发消息
func sendGroupOlineUserMSg(db *gorm.DB, member group_models.GroupMemberModel, msg ctype.Msg, msgID uint) {
	// 查在线的用户列表
	userOnlineIDList := getOnlineUserIDList()
	// 查这个群的成员 并且在线
	var groupMemberOnlineIDList []uint
	db.Model(group_models.GroupMemberModel{}).
		Where("group_id = ? and user_id in ?",
			member.GroupID, userOnlineIDList).
		Select("user_id").Scan(&groupMemberOnlineIDList)
	// 构造响应
	var chatResponse = ChatResponse{
		UserID:         member.UserID,
		Msg:            msg,
		ID:             msgID,
		MsgTyp:         msg.Type,
		CreatedAt:      time.Now(),
		MemberNickname: member.MemberNickname,
	}

	wsInfo, ok := UserOnlineWsMap[member.UserID]
	if ok {
		chatResponse.UserNickname = wsInfo.UserInfo.NickName
		chatResponse.UserAvatar = wsInfo.UserInfo.Avatar
	}
	//向群组在线成员发送消息
	for _, u := range groupMemberOnlineIDList {
		wsUserInfo, ok2 := UserOnlineWsMap[u]
		if !ok2 {
			continue
		}
		chatResponse.IsMe = false
		// 判断isMe
		if wsUserInfo.UserInfo.ID == member.UserID {
			chatResponse.IsMe = true
		}
		byteData, _ := json.Marshal(chatResponse)
		for _, w2 := range wsUserInfo.WsClientMap {
			w2.WriteMessage(websocket.TextMessage, byteData)
		}
	}
}

// 提取在线用户id
func getOnlineUserIDList() (userOnlineIDList []uint) {
	for u, _ := range UserOnlineWsMap {
		userOnlineIDList = append(userOnlineIDList, u)
	}
	return
}

// SendTipErrMsg 发送错误提示消息
func SendTipErrMsg(conn *websocket.Conn, msg string) {

	resp := ChatResponse{
		Msg: ctype.Msg{
			Type: ctype.TipMsgType,
			TipMsg: &ctype.TipMsg{
				Status:  "error",
				Content: msg,
			},
		},
		CreatedAt: time.Now(),
	}
	byteData, _ := json.Marshal(resp)
	conn.WriteMessage(websocket.TextMessage, byteData)
}
