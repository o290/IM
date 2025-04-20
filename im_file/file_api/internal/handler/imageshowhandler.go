package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_file/file_api/internal/logic"
	"server/im_file/file_api/internal/svc"
	"server/im_file/file_api/internal/types"
)

func ImageShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ImageShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewImageShowLogic(r.Context(), svcCtx)
		err := l.ImageShow(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
