package handler

import (
	"net/http"
	"server/common/response"

	"server/im_auth/auth_api/internal/logic"
	"server/im_auth/auth_api/internal/svc"
)

func logoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewLogoutLogic(r.Context(), svcCtx)
		//1.获取token
		token := r.Header.Get("token")
		resp, err := l.Logout(token)
		response.Response(r, w, resp, err)
	}
}
