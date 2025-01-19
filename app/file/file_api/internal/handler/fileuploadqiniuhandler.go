package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"beaver/app/file/file_api/internal/logic"
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"
	"beaver/utils"
	"beaver/utils/md5"
)

// FileTypeMapper maps file extensions to file types.
var FileTypeMapper = map[string]string{
	"jpg":  "image",
	"jpeg": "image",
	"png":  "image",
	"gif":  "image",
	"bmp":  "image",
	"mp4":  "video",
	"avi":  "video",
	"mkv":  "video",
	"mp3":  "audio",
	"wav":  "audio",
	"ogg":  "audio",
	"zip":  "archive",
	"rar":  "archive",
	"7z":   "archive",
	"html": "document",
	"pdf":  "document",
	"doc":  "document",
	"docx": "document",
}

func getFileType(suffix string) string {
	if fileType, ok := FileTypeMapper[suffix]; ok {
		return fileType
	}
	return "unknown"
}

func FileUploadQiniuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		file, fileHead, err := r.FormFile("file")
		if err != nil {
			logx.Error(err)
			response.Response(r, w, nil, err)
			return
		}

		// 文件后缀白名单
		fileName := fileHead.Filename
		nameList := strings.Split(fileName, ".")
		if len(nameList) < 2 {
			response.Response(r, w, nil, errors.New("文件格式不正确"))
			return
		}
		suffix := strings.ToLower(nameList[len(nameList)-1])
		if !utils.InList(svcCtx.Config.WhiteList, suffix) {
			response.Response(r, w, nil, errors.New("文件非法"))
			return
		}

		// 确定文件类型
		fileType := getFileType(suffix)
		if fileType == "unknown" {
			response.Response(r, w, nil, errors.New("未知文件类型"))
			return
		}

		// 检查文件大小
		maxSize, ok := svcCtx.Config.FileMaxSize[fileType]
		if !ok {
			response.Response(r, w, nil, errors.New("配置中未找到该文件类型的最大大小"))
			return
		}
		fileSizeMB := float64(fileHead.Size) / (1024 * 1024)
		if fileSizeMB > maxSize {
			response.Response(r, w, nil, fmt.Errorf("文件大小超过最大限制: %.2fMB", maxSize))
			return
		}

		// 读取文件内容
		byteData, err := io.ReadAll(file)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}
		fileMd5 := md5.MD5(byteData)
		fileMd5Name := fileMd5 + "." + suffix

		l := logic.NewFileUploadQiniuLogic(r.Context(), svcCtx)
		resp, _ := l.FileUploadQiniu(&req)

		host := r.Context().Value("ClientHost")
		scheme := r.Context().Value("Scheme")

		fmt.Println("当前请求的host", host, scheme)

		// 检查文件是否已经存在于数据库中
		var fileModel file_models.FileModel
		err = svcCtx.DB.Take(&fileModel, "hash = ?", fileMd5).Error

		if err == nil {
			resp.FileID = fileModel.FileID
			response.Response(r, w, resp, nil)
			return
		}

		// 根据文件类型创建目录结构，并生成七牛云文件路径
		qiniuFilePath := fmt.Sprintf("%s/%s", fileType, fileMd5Name)

		// 上传文件到七牛云
		qiniuURL, err := uploadToQiniu(qiniuFilePath, byteData, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 创建新的文件记录
		newFileModel := &file_models.FileModel{
			FileName: fileName,
			Size:     fileHead.Size,
			Path:     qiniuURL,
			Hash:     fileMd5,
			FileID:   uuid.New().String(),
			Type:     fileType,
		}
		err = svcCtx.DB.Create(newFileModel).Error
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		resp.FileID = newFileModel.FileID
		response.Response(r, w, resp, nil)
	}
}

func uploadToQiniu(filePath string, fileData []byte, config *svc.ServiceContext) (string, error) {
	// 设置认证信息
	mac := credentials.NewCredentials(config.Config.Qiniu.AK, config.Config.Qiniu.SK)

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})

	reader := bytes.NewReader(fileData)
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: config.Config.Qiniu.Bucket,
		FileName:   filePath,
		ObjectName: &filePath,
	}, nil)

	if err != nil {
		return "", fmt.Errorf("failed to upload file to Qiniu: %v", err)
	}

	return filePath, nil
}
