package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_group/group_api/internal/logic"
	"server/im_group/group_api/internal/svc"
	"server/im_group/group_api/internal/types"
)

func groupSessionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupSessionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGroupSessionLogic(r.Context(), svcCtx)
		resp, err := l.GroupSession(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
