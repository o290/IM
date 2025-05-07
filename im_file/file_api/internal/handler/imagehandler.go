package handler

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"os"
	"path"
	"server/common/response"
	"server/im_file/file_api/internal/logic"
	"server/im_file/file_model"
	"server/utils"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"server/im_file/file_api/internal/svc"
	"server/im_file/file_api/internal/types"
)

// ImageHandler 实现图片上传
func ImageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//1.解析请求参数
		var req types.ImageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		imageType := r.FormValue("imageType")
		switch imageType {
		case "avatar", "group_avatar", "chat": //头像，群头像，聊天头像
		default:
			response.Response(r, w, nil, errors.New("imageType只能为avatar, group_avatar, chat"))
			return
		}

		//2.获取上传的文件图片
		file, fileHead, err := r.FormFile("image")
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		//3.检查文件大小
		mSize := float64(fileHead.Size) / float64(1024) / float64(1024)
		if mSize > svcCtx.Config.FileSize {
			response.Response(r, w, nil, fmt.Errorf("图片大小超过限制,最大只能上传%.2fMB大小的图片", svcCtx.Config.FileSize))
			return
		}
		//4.检查文件后缀是否在白名单
		nameList := strings.Split(fileHead.Filename, ".")
		var suffix string
		if len(nameList) > 1 {
			suffix = nameList[len(nameList)-1]
		}
		if !utils.InitList(svcCtx.Config.WhiteList, suffix) {
			response.Response(r, w, nil, errors.New("图片非法"))
			return
		}

		//5.计算文件哈希值并检查重复
		//先算hash
		imageData, _ := io.ReadAll(file)
		imageHash := utils.MD5(imageData)
		l := logic.NewImageLogic(r.Context(), svcCtx)
		resp, err := l.Image(&req)
		var fileModel file_model.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", imageHash).Error

		if err == nil {
			//找到了,有hash一模一样的,返回之间的那个文件hash组成的web路径
			resp.Url = "/" + fileModel.WebPath()
			logx.Infof("文件%s hash重复", fileHead.Filename)
			return
		}
		//6.创建文件存储目录
		//拼接路径 /uploads/imageType/{uid}.{后缀}
		dirPath := path.Join(svcCtx.Config.UploadDir, imageType)
		_, err = os.ReadDir(dirPath)
		if err != nil {
			os.MkdirAll(dirPath, 0666)
		}

		//7.创建文件结构，并存储文件
		fileName := fileHead.Filename
		newFileModel := file_model.FileModel{
			UserID:   req.UserID,
			FileName: fileName,
			Size:     fileHead.Size,
			Hash:     utils.MD5(imageData),
			Uid:      uuid.New(),
		}
		newFileModel.Path = path.Join(dirPath, fmt.Sprintf("%s.%s", newFileModel.Uid, suffix))
		err = os.WriteFile(newFileModel.Path, imageData, 0666)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		//文件信息入库
		err = svcCtx.DB.Create(&newFileModel).Error
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}
		resp.Url = "/" + newFileModel.WebPath()
		response.Response(r, w, resp, err)
	}
}
