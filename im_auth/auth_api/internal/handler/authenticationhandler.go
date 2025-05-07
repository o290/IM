package handler

import (
	"net/http"
	"server/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_auth/auth_api/internal/logic"
	"server/im_auth/auth_api/internal/svc"
	"server/im_auth/auth_api/internal/types"
)

func authenticationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AuthenticationRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewAuthenticationLogic(r.Context(), svcCtx)
		resp, err := l.Authentication(&req)
		response.Response(r, w, resp, err)
	}
}
