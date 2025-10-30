package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/file"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"
	"beaver/utils"
	"beaver/utils/md5"
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
)

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
	"apk":  "apk",
}

func getFileType(suffix string) string {
	if fileType, ok := FileTypeMapper[suffix]; ok {
		return fileType
	}
	return "unknown"
}

func FileUploadQiniuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("开始处理文件上传请求")

		var req types.FileReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.Error("解析请求参数失败:", err)
			response.Response(r, w, nil, errors.New("解析请求参数失败"))
			return
		}

		file, fileHead, err := r.FormFile("file")
		if err != nil {
			logx.Error("获取上传文件失败:", err)
			response.Response(r, w, nil, errors.New("获取上传文件失败"))
			return
		}
		logx.Info("成功获取上传文件:", fileHead.Filename, "大小:", fileHead.Size)

		// 文件后缀白名单
		originalName := fileHead.Filename
		nameList := strings.Split(originalName, ".")
		if len(nameList) < 2 {
			logx.Error("文件名格式不正确:", originalName)
			response.Response(r, w, nil, errors.New("文件格式不正确"))
			return
		}
		suffix := strings.ToLower(nameList[len(nameList)-1])
		if !utils.InList(svcCtx.Config.WhiteList, suffix) {
			logx.Error("文件类型不在白名单中:", suffix)
			response.Response(r, w, nil, errors.New("文件类型不支持"))
			return
		}
		logx.Info("文件类型检查通过:", suffix)

		// 确定文件类型
		fileType := getFileType(suffix)
		if fileType == "unknown" {
			logx.Error("未知的文件类型:", suffix)
			response.Response(r, w, nil, errors.New("不支持的文件类型"))
			return
		}
		logx.Info("文件类型:", fileType)

		// 检查文件大小
		maxSize, ok := svcCtx.Config.FileMaxSize[fileType]
		if !ok {
			logx.Error("配置中未找到文件类型的大小限制:", fileType)
			response.Response(r, w, nil, errors.New("系统配置错误"))
			return
		}
		fileSizeMB := float64(fileHead.Size) / (1024 * 1024)
		if fileSizeMB > maxSize {
			logx.Error("文件大小超过限制:", fileSizeMB, "MB, 最大限制:", maxSize, "MB")
			response.Response(r, w, nil, fmt.Errorf("文件大小不能超过 %.2fMB", maxSize))
			return
		}
		logx.Info("文件大小检查通过:", fileSizeMB, "MB")

		// 读取文件内容
		logx.Info("开始读取文件内容")
		byteData, err := io.ReadAll(file)
		if err != nil {
			logx.Error("读取文件内容失败:", err)
			response.Response(r, w, nil, errors.New("读取文件失败"))
			return
		}
		logx.Info("文件内容读取成功, 大小:", len(byteData), "字节")

		fileMd5 := md5.MD5(byteData)
		fileMd5Name := fileMd5 + "." + suffix
		logx.Info("文件MD5:", fileMd5)

		l := logic.NewFileUploadQiniuLogic(r.Context(), svcCtx)
		resp, _ := l.FileUploadQiniu(&req)

		// 检查文件是否已经存在于数据库中
		logx.Info("检查文件是否已存在")
		var fileModel file_models.FileModel
		err = svcCtx.DB.Take(&fileModel, "md5 = ?", fileMd5).Error

		if err == nil {
			logx.Info("文件已存在，直接返回:", fileModel.FileName, fileModel.OriginalName)
			resp.FileName = fileModel.FileName
			resp.OriginalName = fileModel.OriginalName
			response.Response(r, w, resp, nil)
			return
		}
		logx.Info("文件不存在，继续上传流程")

		// 根据文件类型创建目录结构，并生成七牛云文件路径
		qiniuFilePath := fmt.Sprintf("%s/%s", fileType, fileMd5Name)
		logx.Info("七牛云文件路径:", qiniuFilePath)

		// 上传文件到七牛云
		logx.Info("开始上传文件到七牛云")
		qiniuURL, err := uploadToQiniu(qiniuFilePath, byteData, svcCtx)
		if err != nil {
			logx.Error("上传到七牛云失败:", err)
			response.Response(r, w, nil, errors.New("上传文件失败"))
			return
		}
		logx.Info("文件成功上传到七牛云")

		// 创建新的文件记录
		logx.Info("开始创建数据库记录")
		// 生成带后缀的FileName
		fileUUID := uuid.New().String()
		fileNameWithSuffix := fileUUID + "." + suffix

		newFileModel := &file_models.FileModel{
			OriginalName: strings.TrimSuffix(originalName, "."+suffix),
			Size:         fileHead.Size,
			Path:         qiniuURL,
			Md5:          fileMd5,
			FileName:     fileNameWithSuffix,
			Type:         fileType,
		}
		err = svcCtx.DB.Create(newFileModel).Error
		if err != nil {
			logx.Error("创建数据库记录失败:", err)
			response.Response(r, w, nil, errors.New("保存文件信息失败"))
			return
		}
		logx.Info("数据库记录创建成功:", newFileModel.FileName)

		resp.FileName = newFileModel.FileName
		resp.OriginalName = newFileModel.OriginalName

		logx.Info("文件上传完成:", resp.FileName)
		response.Response(r, w, resp, nil)
	}
}

func uploadToQiniu(filePath string, fileData []byte, config *svc.ServiceContext) (string, error) {
	logx.Info("准备上传到七牛云, 文件路径:", filePath)

	// 设置认证信息
	mac := credentials.NewCredentials(config.Config.Qiniu.AK, config.Config.Qiniu.SK)
	logx.Info("七牛云认证信息设置完成")

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	logx.Info("七牛云上传管理器创建完成")

	reader := bytes.NewReader(fileData)
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: config.Config.Qiniu.Bucket,
		FileName:   filePath,
		ObjectName: &filePath,
	}, nil)

	if err != nil {
		logx.Error("七牛云上传失败:", err)
		return "", fmt.Errorf("failed to upload file to Qiniu: %v", err)
	}
	logx.Info("七牛云上传成功")

	return filePath, nil
}
