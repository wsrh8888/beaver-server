package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"beaver/app/file/file_api/internal/handler/common"
	"beaver/app/file/file_api/internal/logic"
	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"beaver/common/response"
)

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

		// 使用公共工具函数验证和处理文件
		fileReq, err := common.ValidateAndProcessFile(file, fileHead, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFileUploadQiniuLogic(r.Context(), svcCtx)
		resp, _ := l.FileUploadQiniu(&req)

		// 检查文件是否已经存在于数据库中
		existingFile, err := common.CheckFileExists(fileReq.FileMd5, svcCtx)
		if err == nil {
			resp.FileName = existingFile.FileName
			resp.OriginalName = existingFile.OriginalName

			// 如果文件已存在但FileInfo为空，可以设置默认值
			if existingFile.FileInfo == nil {
				existingFile.FileInfo = &file_models.FileInfo{
					Type: file_models.FileType(existingFile.Type),
				}
				svcCtx.DB.Save(existingFile)
				logx.Infof("已存在文件元数据设置默认值: %s", existingFile.FileName)
			}

			// 转换FileInfo为API响应格式
			if existingFile.FileInfo != nil {
				resp.FileInfo = common.ConvertFileInfoToAPI(existingFile.FileInfo)
			}

			response.Response(r, w, resp, nil)
			return
		}

		// 根据文件类型创建目录结构，并生成七牛云文件路径
		fileMd5Name := fileReq.FileMd5 + "." + fileReq.Suffix
		qiniuFilePath := fmt.Sprintf("%s/%s", fileReq.FileType, fileMd5Name)

		// 上传文件到七牛云
		qiniuURL, err := uploadToQiniu(qiniuFilePath, fileReq.ByteData, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 创建文件记录
		newFileModel, err := common.CreateFileRecord(fileReq, qiniuURL, file_models.QiniuSource, svcCtx)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 初始化文件信息
		var fileInfo *file_models.FileInfo

		// 手动解析FormData中的fileInfo字段
		if fileInfoStr := r.FormValue("fileInfo"); fileInfoStr != "" {
			var apiFileInfo types.FileInfo
			if err := json.Unmarshal([]byte(fileInfoStr), &apiFileInfo); err == nil {
				fileInfo = common.ConvertAPIFileInfoToModel(&apiFileInfo)
			}
		}

		if fileInfo != nil {
			// 更新数据库中的FileInfo
			newFileModel.FileInfo = fileInfo
			svcCtx.DB.Save(newFileModel)
			logx.Infof("文件元数据获取成功: %s", newFileModel.FileName)
		}

		resp.FileName = newFileModel.FileName
		resp.OriginalName = newFileModel.OriginalName

		// 转换FileInfo为API响应格式
		if fileInfo != nil {
			resp.FileInfo = common.ConvertFileInfoToAPI(fileInfo)
		}

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
