package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_user/user_api/internal/logic"
	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"
)

func userValidListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FriendValidRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserValidListLogic(r.Context(), svcCtx)
		resp, err := l.UserValidList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
