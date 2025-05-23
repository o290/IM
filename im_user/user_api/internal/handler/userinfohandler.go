package handler

import (
	"net/http"
	"server/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_user/user_api/internal/logic"
	"server/im_user/user_api/internal/svc"
	"server/im_user/user_api/internal/types"
)

func UserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.UserInfo(&req)
		response.Response(r, w, resp, err)
	}
}
