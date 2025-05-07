package handler

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"server/common/response"

	"server/im_auth/auth_api/internal/logic"
	"server/im_auth/auth_api/internal/svc"
	"server/im_auth/auth_api/internal/types"
)

func loginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRequest

		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		l := logic.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		response.Response(r, w, resp, err)
	}
}
