// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3

package handler

import (
	"net/http"

	"server/im_chat/chat_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodDelete,
				Path:    "/api/chat/chat",
				Handler: ChatDeleteHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/chat/history",
				Handler: ChatHistoryHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/chat/session",
				Handler: ChatSessionHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/user_top",
				Handler: UserTopHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/chat/ws/chat",
				Handler: ChatHandler(serverCtx),
			},
		},
	)
}
