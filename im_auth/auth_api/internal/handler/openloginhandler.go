package handler

import (
	"net/http"
	"server/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_auth/auth_api/internal/logic"
	"server/im_auth/auth_api/internal/svc"
	"server/im_auth/auth_api/internal/types"
)

func open_loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenLoginRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOpen_loginLogic(r.Context(), svcCtx)
		resp, err := l.Open_login(&req)
		//if err != nil {
		//	httpx.ErrorCtx(r.Context(), w, err)
		//} else {
		//	httpx.OkJsonCtx(r.Context(), w, resp)
		//}
		response.Response(r, w, resp, err)
	}
}
