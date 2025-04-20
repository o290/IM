package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
	"net/http"
	"server/common/models/ctype"
	"server/common/response"
	"server/common/service/redis_service"
	"server/im_chat/chat_models"
	"server/im_file/file_rpc/files"
	"server/im_user/user_models"
	"server/im_user/user_rpc/types/user_rpc"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_chat/chat_api/internal/svc"
	"server/im_chat/chat_api/internal/types"
)

// UserWsInfo 用户websocket信息
type UserWsInfo struct {
	UserInfo      user_models.UserModel      //用户信息
	WsClientMap   map[string]*websocket.Conn //这个用户管理的所有客户端
	CurrentConn   *websocket.Conn            //当前的连接对象
	LastHeartbeat time.Time                  // 最后一次收到心跳的时间
}

var UserOnlineWsMap = map[uint]*UserWsInfo{}

const HeartbeatInterval = 10 * time.Second // 心跳间隔
const HeartbeatTimeout = 30 * time.Second  // 心跳超时时间

// ChatHandler 聊天处理
func ChatHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 一.初始化与处理请求
		// 1.从header中解析请求，得到请求体req
		var req types.ChatRequest
		// 解析handler
		if err := httpx.ParseHeaders(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 2.创建一个 websocket.Upgrader 实例用于将 HTTP 连接升级为 WebSocket 连接
		var upGrader = websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// CheckOrigin 返回 true，允许将该 HTTP 连接升级为 WebSocket 连接
				// CheckOrigin 则是用来对跨域请求进行控制
				return true
			},
		}

		// 3.将http升级为websocket
		conn, err := upGrader.Upgrade(w, r, nil)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		// 二.用户连接管理
		// 1.RemoteAddr获取当前连接的远程地址，并转换字符串
		addr := conn.RemoteAddr().String()
		// 2.处理websocket连接关闭
		defer func() {
			conn.Close()
			// 获取当前用户的在线信息，在线则退出连接
			userWsInfo, ok := UserOnlineWsMap[req.UserID]
			if ok {
				// 删除退出的ws信息
				delete(userWsInfo.WsClientMap, addr)
			}
			// 检查该用户的所有连接是否全部关闭
			if userWsInfo != nil && len(userWsInfo.WsClientMap) == 0 {
				// 全退完了，删除对应缓存
				delete(UserOnlineWsMap, req.UserID)
				svcCtx.Redis.HDel(context.Background(), "online", fmt.Sprintf("%d", req.UserID))
			}
		}()
		// 3.调用户服务，获取当前用户信息
		res, err := svcCtx.UserRpc.UserInfo(context.Background(), &user_rpc.UserInfoRequest{
			UserId: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		// 将服务返回的 JSON 数据解析到 userInfo 结构体中
		var userInfo user_models.UserModel
		json.Unmarshal(res.Data, &userInfo)
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		// 4.更新用户的在线信息
		// 查找用户是否在线，若不在线，表示这是用户的首次连接
		userWsInfo, ok := UserOnlineWsMap[req.UserID]
		if !ok {
			// 用户首次连接，存入map中
			userWsInfo = &UserWsInfo{
				UserInfo:      userInfo,
				WsClientMap:   map[string]*websocket.Conn{addr: conn},
				CurrentConn:   conn,       // 当前连接对象
				LastHeartbeat: time.Now(), // 记录首次连接时间
			}
			UserOnlineWsMap[req.UserID] = userWsInfo
		}
		// 检查当前客户端地址是否已经存在于用户的 WsClientMap 中
		_, ok1 := userWsInfo.WsClientMap[addr]
		// 如果不存在，说明用户是第二次或更多次连接
		if !ok1 {
			UserOnlineWsMap[req.UserID].WsClientMap[addr] = conn
			UserOnlineWsMap[req.UserID].CurrentConn = conn
		}
		// 把在线用户存进redis
		svcCtx.Redis.HSet(context.Background(), "online", fmt.Sprintf("%d", req.UserID), req.UserID)
		// 遍历在线的用户，和当前这个人是好友的，就给他发送好友在线
		// 先取出所有在线的用户id取出来,以及待确认的用户id,然后传到用户rpc服务中
		// 在rpc服务中,去判断哪些用户是好友关系

		// 三.在用户上线时，获取该用户的好友列表，检查好友是否在线以及是否开启了好友上线提醒功能，
		// 如果满足条件则向在线且开启提醒的好友发送上线通知
		// 获取好友列表
		friendRes, err := svcCtx.UserRpc.FriendList(context.Background(), &user_rpc.FriendListRequest{
			User: uint32(req.UserID),
		})
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		logx.Infof("用户上线：%s ,用户id：%d", userInfo.Nickname, req.UserID)
		// 遍历好友列表
		for _, info := range friendRes.FriendList {
			// 查看好友的在线信息，若在线，构造一条上线通知消息，包含当前上线用户的昵称
			friend, ok := UserOnlineWsMap[uint(info.UserId)]
			if ok {
				text := fmt.Sprintf("好友%s上线了", UserOnlineWsMap[req.UserID].UserInfo.Nickname)
				// 判断用户是否开了好友上线提醒
				if friend.UserInfo.UserConfModel.FriendOnline {
					// 利用好友的websocket客户端给自己发送上线通知
					sendWsMapMsg(friend.WsClientMap, []byte(text))
				}
			}
		}

		// 启动心跳检测协程
		go heartbeatCheck(svcCtx, req.UserID, userWsInfo)

		// 四.循环从websocket连接conn中读取消息，消息预处理，根据不同类型消息处理，然后存储并发送消息
		for {
			// 读取消息：消息类型 消息 错误
			messageType, p, err1 := conn.ReadMessage()
			if err1 != nil {
				// 用户断开聊天
				break
			}
			if messageType == websocket.PingMessage {
				// 收到心跳消息，更新最后心跳时间
				userWsInfo.LastHeartbeat = time.Now()
				// 回复心跳消息
				err := conn.WriteMessage(websocket.PongMessage, []byte{})
				if err != nil {
					logx.Error("Failed to send heartbeat response:", err)
					break
				}
				continue
			}

			// 解析消息
			var request ChatRequest
			err2 := json.Unmarshal(p, &request)
			if err2 != nil {
				// 用户乱发消息
				logx.Error(err2)
				SendTipErrMsg(conn, "参数解析失败")
				continue
			}

			// 检验聊天对象是否是好友
			// 消息的接收用户 ID request.RevUserID 不等于当前用户 ID req.UserID
			if request.RevUserID != req.UserID {
				// 判断聊天的是否是你的好友
				isFriendRes, err := svcCtx.UserRpc.IsFriend(context.Background(), &user_rpc.IsFriendRequest{
					User1: uint32(req.UserID),
					User2: uint32(request.RevUserID),
				})
				if err != nil {
					logx.Error(err2)
					SendTipErrMsg(conn, "用户服务错误")
					continue
				}
				if !isFriendRes.IsFriend {
					SendTipErrMsg(conn, "你们还不是好友呢")
					continue
				}
			}

			// 判断type
			if !(request.Msg.Type >= 1 && request.Msg.Type <= 12) {
				SendTipErrMsg(conn, "消息类型错误")
				continue
			}
			// 根据消息类型进行处理
			switch request.Msg.Type {
			case ctype.TextMsgType:
				// 消息内容缺失
				if request.Msg.TextMsg == nil {
					SendTipErrMsg(conn, "请输入消息内容")
					continue
				}
				// 消息内容为空
				if request.Msg.TextMsg.Content == "" {
					SendTipErrMsg(conn, "请输入消息内容")
					continue
				}
			case ctype.FileMsgType:
				if request.Msg.FileMsg == nil {
					SendTipErrMsg(conn, "请上传文件")
					return
				}
				// 如果是文件类型,那么就要去请求文件rpc服务
				nameList := strings.Split(request.Msg.FileMsg.Src, "/")
				if len(nameList) == 0 {
					SendTipErrMsg(conn, "请上传文件")
					continue
				}
				// 文件名fileID
				fileID := nameList[len(nameList)-1]
				// 获取文件信息
				fileResponse, err3 := svcCtx.FileRpc.FileInfo(context.Background(), &files.FileInfoRequest{
					FileId: fileID,
				})
				if err3 != nil {
					logx.Error(err3)
					SendTipErrMsg(conn, err3.Error())
					continue
				}
				request.Msg.FileMsg.Title = fileResponse.FileName
				request.Msg.FileMsg.Size = fileResponse.FileSize
				request.Msg.FileMsg.Type = fileResponse.FileType
			case ctype.WithdrawMsgType:
				// 撤回消息的消息id是必填的
				if request.Msg.WithdrawMsg.MsgID == 0 {
					SendTipErrMsg(conn, "撤回消息id必填")
					continue
				}
				// 自己只能撤回自己的
				//  找这个消息是谁发的
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.WithdrawMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}
				// 已经是撤回消息的,不能再撤回了
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "撤回消息不能再撤回了")
					continue
				}
				// 判断是不是自己发的
				if msgModel.SendUserID != req.UserID {
					SendTipErrMsg(conn, "只能撤回自己的消息")
					continue
				}

				// 判断消息的时间,小于2分钟的才能撤回
				now := time.Now()
				subTime := now.Sub(msgModel.CreatedAt)
				if subTime >= time.Minute*2 {
					SendTipErrMsg(conn, "只能撤回两分钟以内的消息哦")
					continue
				}
				// 撤回逻辑
				// 收到撤回请求之后,服务端这边把原消息类型修改为撤回消息类型,并且记录原消息
				//  然后通知前端的收发双方,重新拉取聊天记录
				var content = "撤回了一条消息"
				if userInfo.UserConfModel.RecallMessage != nil {
					content = *userInfo.UserConfModel.RecallMessage
				}
				content = "你" + content
				// 前端可以判断,这个消息如果不是isMe,就可以把你替换成对方的昵称
				originMsg := msgModel.Msg
				// 这里可能会出现循环引用,所以拷贝了这个值,并且把撤回消息置空了
				originMsg.WithdrawMsg = nil
				svcCtx.DB.Model(&msgModel).Updates(chat_models.ChatModel{
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
				// 回复消息
				// 先校验
				if request.Msg.ReplyMsg == nil || request.Msg.ReplyMsg.MsgID == 0 {
					SendTipErrMsg(conn, "回复消息id必填")
					return
				}
				//  找这个原消息
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.ReplyMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}

				// 不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}
				// 回复安全性问题
				// 回复的这个消息,必须是你自己或者当前和你聊天这个人发出来的

				// 原消息必须是 当前你要和对方聊的 原消息就会有一个 发送人id和接收入id,我们聊天也会有一个发送人id和接收人id
				// 因为回复消息可以回复自己的,也可以回复别人的
				// 如果回复只能回复别人的?那么条件怎么写?
				if !((msgModel.SendUserID == req.UserID && msgModel.RevUserID == request.RevUserID) ||
					(msgModel.SendUserID == request.RevUserID && msgModel.RevUserID == req.UserID)) {
					SendTipErrMsg(conn, "只能回复自己或者对方的消息")
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
				// 回复消息
				// 先校验
				if request.Msg.QuoteMsg == nil || request.Msg.QuoteMsg.MsgID == 0 {
					SendTipErrMsg(conn, "引用消息id必填")
					return
				}
				//  找这个原消息
				var msgModel chat_models.ChatModel
				err = svcCtx.DB.Take(&msgModel, request.Msg.QuoteMsg.MsgID).Error
				if err != nil {
					SendTipErrMsg(conn, "消息不存在")
					continue
				}
				// 不能回复撤回消息
				if msgModel.MsgType == ctype.WithdrawMsgType {
					SendTipErrMsg(conn, "该消息已撤回")
					continue
				}
				// 回复安全性问题
				// 回复的这个消息,必须是你自己或者当前和你聊天这个人发出来的

				// 原消息必须是 当前你要和对方聊的 原消息就会有一个 发送人id和接收入id,我们聊天也会有一个发送人id和接收人id
				// 因为回复消息可以回复自己的,也可以回复别人的
				// 如果回复只能回复别人的?那么条件怎么写?
				if !((msgModel.SendUserID == req.UserID && msgModel.RevUserID == request.RevUserID) ||
					(msgModel.SendUserID == request.RevUserID && msgModel.RevUserID == req.UserID)) {
					SendTipErrMsg(conn, "只能回复自己或者对方的消息")
					continue
				}
				userBaseInfo, err5 := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, msgModel.SendUserID)

				if err5 != nil {
					logx.Error(err)
					SendTipErrMsg(conn, err5.Error())
					continue
				}
				request.Msg.QuoteMsg.Msg = &msgModel.Msg
				request.Msg.QuoteMsg.UserID = msgModel.SendUserID
				request.Msg.QuoteMsg.UserNickName = userBaseInfo.NickName
				request.Msg.QuoteMsg.OriginMsgDate = msgModel.CreatedAt
			}
			// 消息入库,入库就是会把聊天记录保存到数据库中
			msgID := InsertMsgByChat(svcCtx.DB, request.RevUserID, req.UserID, request.Msg)
			// 将消息发送给发送者和接受者，看看目标用户在不在线 给发送双方都要发送消息
			SendMsgByUser(svcCtx, request.RevUserID, req.UserID, request.Msg, msgID)
		}
	}
}

// heartbeatCheck 心跳检测协程
func heartbeatCheck(svcCtx *svc.ServiceContext, userId uint, userWsInfo *UserWsInfo) {
	ticker := time.NewTicker(HeartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		if time.Since(userWsInfo.LastHeartbeat) > HeartbeatTimeout {
			// 心跳超时，认为用户离线
			for addr := range userWsInfo.WsClientMap {
				userWsInfo.WsClientMap[addr].Close()
				delete(userWsInfo.WsClientMap, addr)
			}
			if len(userWsInfo.WsClientMap) == 0 {
				delete(UserOnlineWsMap, userId)
				svcCtx.Redis.HDel(context.Background(), "online", fmt.Sprintf("%d", userId))
			}
			logx.Infof("用户离线：%s ,用户id：%d", userWsInfo.UserInfo.Nickname, userId)
			break
		}
	}
}

type ChatRequest struct {
	RevUserID uint      `json:"revUserID"` //给谁发
	Msg       ctype.Msg `json:"msg"`
}
type ChatResponse struct {
	ID        uint           `json:"id"`
	IsMe      bool           `json:"isMe"`
	RevUser   ctype.UserInfo `json:"revUser"`
	SendUser  ctype.UserInfo `json:"sendUser"`
	Msg       ctype.Msg      `json:"msg"`
	CreatedAt time.Time      `json:"created_at"`
}

// sendWsMapMsg 向一个存储了多个 WebSocket 连接的映射中的所有连接发送相同的文本消息
func sendWsMapMsg(wsMap map[string]*websocket.Conn, byteData []byte) {
	for _, conn := range wsMap {
		conn.WriteMessage(websocket.TextMessage, byteData)
	}
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

// InsertMsgByChat 消息入库
func InsertMsgByChat(db *gorm.DB, revUserId uint, sendUserID uint, msg ctype.Msg) (msgID uint) {
	switch msg.Type {
	// 处理撤回消息
	case ctype.WithdrawMsgType:
		fmt.Println("撒回消息自己是不入库的")
		return
	}
	chatModel := chat_models.ChatModel{
		SendUserID: sendUserID,
		RevUserID:  revUserId,
		MsgType:    msg.Type,
		Msg:        msg,
	}
	// 生成消息预览内容
	chatModel.MsgPreview = chatModel.MsgPreviewMethod()
	// 创建消息
	err := db.Create(&chatModel).Error
	if err != nil {
		logx.Error(err)
		sendUser, ok := UserOnlineWsMap[sendUserID]
		if !ok {
			return
		}
		SendTipErrMsg(sendUser.CurrentConn, "消息保存失败")
	}
	return chatModel.ID
}

// SendMsgByUser 处理用户之间的消息发送,根据消息的发送者和接收者的在线状态，将消息发送给相应的用户
func SendMsgByUser(svcCtx *svc.ServiceContext, revUserId uint, sendUserID uint, msg ctype.Msg, msgID uint) {
	// 检查用户在线状态，获取接受者和发送者的消息
	revUser, ok1 := UserOnlineWsMap[revUserId]
	sendUser, ok2 := UserOnlineWsMap[sendUserID]
	// 构造响应
	resp := ChatResponse{
		ID:        msgID,
		Msg:       msg,
		CreatedAt: time.Now(),
	}
	// 自己给自己发送消息的情况
	if ok1 && ok2 && sendUserID == revUserId {
		// 自己给自己发
		resp.RevUser = ctype.UserInfo{
			ID:       revUserId,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		}
		resp.SendUser = ctype.UserInfo{
			ID:       sendUserID,
			NickName: sendUser.UserInfo.Nickname,
			Avatar:   sendUser.UserInfo.Avatar,
		}
		byteData, _ := json.Marshal(resp)
		//revUser.Conn.WriteMessage(websocket.TextMessage, byteData)
		sendWsMapMsg(revUser.WsClientMap, byteData)
		return
	}

	if !ok1 {
		// 处理接受者不在的情况
		userBaseInfo, err := redis_service.GetUserBaseInfo(svcCtx.Redis, svcCtx.UserRpc, revUserId)
		if err != nil {
			logx.Error(err)
			return
		}
		// 接受者信息
		resp.RevUser = ctype.UserInfo{
			ID:       revUserId,
			NickName: userBaseInfo.NickName,
			Avatar:   userBaseInfo.Avatar,
		}
	} else {
		// 处理接受者在的情况
		resp.RevUser = ctype.UserInfo{
			ID:       revUserId,
			NickName: revUser.UserInfo.Nickname,
			Avatar:   revUser.UserInfo.Avatar,
		}
	}

	// 给发送者发送消息
	// 发送者在线
	resp.SendUser = ctype.UserInfo{
		ID:       sendUserID,
		NickName: sendUser.UserInfo.Nickname,
		Avatar:   sendUser.UserInfo.Avatar,
	}
	resp.IsMe = true
	byteData, _ := json.Marshal(resp)

	sendWsMapMsg(sendUser.WsClientMap, byteData)
	if ok1 {
		//  接收者在线
		resp.IsMe = false
		byteData, _ = json.Marshal(resp)
		sendWsMapMsg(revUser.WsClientMap, byteData)
	}
}
