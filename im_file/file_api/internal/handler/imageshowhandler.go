package handler

import (
	"errors"
	"net/http"
	"os"
	"server/common/response"
	"server/im_file/file_model"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_file/file_api/internal/svc"
	"server/im_file/file_api/internal/types"
)

func ImageShowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//1.解析请求参数
		var req types.ImageShowRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		//2.从数据库中查找图片信息
		var fileModel file_model.FileModel
		err := svcCtx.DB.Take(&fileModel, "uid = ?", req.ImageName).Error
		if err != nil {
			response.Response(r, w, nil, errors.New("文件不存在"))
			return

		}

		//3.读取图片文件内容
		byteData, err := os.ReadFile(fileModel.Path)
		if err != nil {
			response.Response(r, w, nil, err)
		}
		w.Write(byteData)
	}
}
