package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"os"
	"path"
	"server/common/response"
	"server/im_file/file_model"
	"server/im_user/user_rpc/types/user_rpc"
	"server/utils"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_file/file_api/internal/logic"
	"server/im_file/file_api/internal/svc"
	"server/im_file/file_api/internal/types"
)

// FileHandler 文件上传
func FileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//1.解析请求头
		var req types.FileRequest
		if err := httpx.ParseHeaders(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		//2.获取上传文件
		file, fileHead, err := r.FormFile("file")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		//3.检查文件后缀名是否在黑名单中
		//文件上传用黑名单 exe php
		//文件后缀白名单
		nameList := strings.Split(fileHead.Filename, ".")
		var suffix string
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}
		if utils.InitList(svcCtx.Config.BlackList, suffix) {
			response.Response(r, w, nil, errors.New("文件非法"))
			return
		}

		//4.计算文件哈希值
		fileData, _ := io.ReadAll(file)
		fileHash := utils.MD5(fileData)

		//5.检查文件哈希值是否重复
		l := logic.NewFileLogic(r.Context(), svcCtx)
		resp, err := l.File(&req)

		var fileModel file_model.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", fileHash).Error
		if err == nil {
			resp.Src = fileModel.WebPath()
			logx.Infof("文件%s hash重复", fileHead.Filename)
			response.Response(r, w, resp, err)
			return
		}
		//文件重命名
		//在保存文件之前,去读一些文件列表 如果有重名的,算一下它们两个的hash值,一样的就不用写了
		//他们的hash如果不一样，就把最新的这个重命名为old_name_xxx.xxx

		//6.获取用户信息并创建文件目录
		userResponse, err := svcCtx.UserRpc.UserListInfo(context.Background(), &user_rpc.UserListInfoRequest{
			UserIdList: []uint32{uint32(req.UserID)},
		})
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		dirName := fmt.Sprintf("%d_%s", req.UserID, userResponse.UserInfo[uint32(req.UserID)].NickName)
		dirPath := path.Join(svcCtx.Config.UploadDir, "file", dirName)
		_, err = os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}
		//7.创建文件模型并保存文件
		newFileModel := file_model.FileModel{
			UserID:   req.UserID,
			FileName: fileHead.Filename,
			Size:     fileHead.Size,
			Hash:     fileHash,
			Uid:      uuid.New(),
		}
		newFileModel.Path = path.Join(dirPath, fmt.Sprintf("%s.%s", newFileModel.Uid, suffix))

		//8.将文件信息入库
		err = os.WriteFile(newFileModel.Path, fileData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		err = svcCtx.DB.Create(&newFileModel).Error
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		resp.Src = fileModel.WebPath()
		response.Response(r, w, resp, err)
	}
}
