/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	filecommon "beaver/app/backend/backend_admin/internal/handler/file/common"
	logic "beaver/app/backend/backend_admin/internal/logic/file"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	utils "beaver/utils/list"
	"beaver/utils/md5"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadQiniuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := beaverlog.New("file_upload_qiniu", r.Context())

		logger.Info(model.LogMsg{Text: "开始处理文件上传请求"})

		var req types.FileUploadQiniuReq
		if err := httpx.Parse(r, &req); err != nil {
			logger.Error(model.LogMsg{Text: "解析请求参数失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("解析请求参数失败"))
			return
		}

		file, fileHead, err := r.FormFile("file")
		if err != nil {
			logger.Error(model.LogMsg{Text: "获取上传文件失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("获取上传文件失败"))
			return
		}
		logger.Info(model.LogMsg{Text: "成功获取上传文件", Data: map[string]interface{}{"filename": fileHead.Filename, "size": fileHead.Size}})

		// 获取fileInfo
		fileInfoStr := r.FormValue("fileInfo")

		// 确定文件来源
		source := "qiniu"
		if req.Source != "" && req.Source == "local" {
			source = "local"
		}

		// 文件后缀白名单
		originalName := fileHead.Filename
		nameList := strings.Split(originalName, ".")
		if len(nameList) < 2 {
			logger.Error(model.LogMsg{Text: "文件名格式不正确", Data: map[string]interface{}{"name": originalName}})
			response.Response(r, w, nil, errors.New("文件格式不正确"))
			return
		}
		suffix := strings.ToLower(nameList[len(nameList)-1])
		if !utils.InList(svcCtx.Config.File.WhiteList, suffix) {
			logger.Error(model.LogMsg{Text: "文件类型不在白名单中", Data: map[string]interface{}{"suffix": suffix}})
			response.Response(r, w, nil, errors.New("文件类型不支持"))
			return
		}
		logger.Info(model.LogMsg{Text: "文件类型检查通过", Data: map[string]interface{}{"suffix": suffix}})

		// 确定文件类型
		fileType := filecommon.GetFileType(suffix)
		if fileType == "unknown" {
			logger.Error(model.LogMsg{Text: "未知的文件类型", Data: map[string]interface{}{"suffix": suffix}})
			response.Response(r, w, nil, errors.New("不支持的文件类型"))
			return
		}
		logger.Info(model.LogMsg{Text: "文件类型", Data: map[string]interface{}{"type": fileType}})

		// 检查文件大小
		maxSize, ok := svcCtx.Config.File.MaxSize[fileType]
		if !ok {
			logger.Error(model.LogMsg{Text: "未找到文件类型大小限制", Data: map[string]interface{}{"type": fileType}})
			response.Response(r, w, nil, errors.New("系统配置错误"))
			return
		}
		fileSizeMB := float64(fileHead.Size) / (1024 * 1024)
		if fileSizeMB > maxSize {
			logger.Error(model.LogMsg{Text: "文件大小超过限制", Data: map[string]interface{}{"sizeMB": fileSizeMB, "maxMB": maxSize}})
			response.Response(r, w, nil, fmt.Errorf("文件大小不能超过 %.2fMB", maxSize))
			return
		}
		logger.Info(model.LogMsg{Text: "文件大小检查通过", Data: map[string]interface{}{"sizeMB": fileSizeMB}})

		// 读取文件内容
		logger.Info(model.LogMsg{Text: "开始读取文件内容"})
		byteData, err := io.ReadAll(file)
		if err != nil {
			logger.Error(model.LogMsg{Text: "读取文件内容失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("读取文件失败"))
			return
		}
		logger.Info(model.LogMsg{Text: "文件内容读取成功", Data: map[string]interface{}{"bytes": len(byteData)}})

		fileMd5 := md5.MD5(byteData)
		fileMd5Name := fileMd5 + "." + suffix
		logger.Info(model.LogMsg{Text: "文件MD5", Data: map[string]interface{}{"md5": fileMd5}})

		l := logic.NewFileUploadQiniuLogic(r.Context(), svcCtx)
		resp, err := l.FileUploadQiniu(&req)
		if err != nil {
			logger.Error(model.LogMsg{Text: "七牛文件上传处理失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("上传文件失败"))
			return
		}

		// 检查文件是否已经存在于数据库中
		logger.Info(model.LogMsg{Text: "检查文件是否已存在"})
		existingFile, found, err := filecommon.FindFileByMd5(r.Context(), fileMd5, svcCtx)
		if err != nil {
			logger.Error(model.LogMsg{Text: "查询文件失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("查询文件失败"))
			return
		}
		if found {
			logger.Info(model.LogMsg{Text: "文件已存在直接返回", Data: map[string]interface{}{"fileKey": existingFile.FileKey, "name": existingFile.OriginalName}})
			resp.OriginalName = existingFile.OriginalName
			resp.FileURL = filecommon.BuildQiniuFileURL(svcCtx.Config.Qiniu.Domain, existingFile.Path)
			response.Response(r, w, resp, nil)
			return
		}
		logger.Info(model.LogMsg{Text: "文件不存在继续上传"})

		if fileInfoStr == "" {
			logger.Error(model.LogMsg{Text: "fileInfo不能为空"})
			response.Response(r, w, nil, errors.New("fileInfo不能为空"))
			return
		}
		if !json.Valid([]byte(fileInfoStr)) {
			logger.Error(model.LogMsg{Text: "fileInfo格式不正确"})
			response.Response(r, w, nil, errors.New("fileInfo格式不正确"))
			return
		}

		// 根据文件类型创建目录结构，并生成七牛云文件路径
		// 如果配置了项目名称，则添加项目目录前缀；否则使用根目录
		projectName := svcCtx.Config.Qiniu.ProjectName
		var qiniuFilePath string
		if projectName != "" {
			qiniuFilePath = fmt.Sprintf("%s/%s/%s", projectName, fileType, fileMd5Name)
		} else {
			qiniuFilePath = fmt.Sprintf("%s/%s", fileType, fileMd5Name)
		}
		logger.Info(model.LogMsg{Text: "七牛云文件路径", Data: map[string]interface{}{"path": qiniuFilePath}})

		// 上传文件到七牛云
		logger.Info(model.LogMsg{Text: "开始上传文件到七牛云"})
		qiniuURL, err := uploadToQiniu(qiniuFilePath, byteData, svcCtx, logger)
		if err != nil {
			logger.Error(model.LogMsg{Text: "上传到七牛云失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("上传文件失败"))
			return
		}
		logger.Info(model.LogMsg{Text: "文件成功上传到七牛云"})

		saveReq := &types.SaveFileReq{
			OriginalName: strings.TrimSuffix(originalName, "."+suffix),
			Size:         fileHead.Size,
			Path:         qiniuURL,
			Md5:          fileMd5,
			Type:         fileType,
			Source:       source,
			FileInfo:     fileInfoStr,
		}
		saveLogic := logic.NewSaveFileLogic(r.Context(), svcCtx)
		saveResp, err := saveLogic.SaveFile(saveReq)
		if err != nil {
			logger.Error(model.LogMsg{Text: "保存文件信息失败", Data: map[string]interface{}{"err": err.Error()}})
			response.Response(r, w, nil, errors.New("保存文件信息失败"))
			return
		}
		logger.Info(model.LogMsg{Text: "数据库记录创建成功", Data: map[string]interface{}{"fileKey": saveResp.FileKey}})

		resp.OriginalName = saveReq.OriginalName
		resp.FileURL = filecommon.BuildQiniuFileURL(svcCtx.Config.Qiniu.Domain, qiniuURL)

		logger.Info(model.LogMsg{Text: "文件上传完成", Data: map[string]interface{}{"url": resp.FileURL}})
		response.Response(r, w, resp, nil)
	}
}

func uploadToQiniu(filePath string, fileData []byte, config *svc.ServiceContext, logger *beaverlog.Logger) (string, error) {
	logger.Info(model.LogMsg{Text: "准备上传到七牛云", Data: map[string]interface{}{"path": filePath}})

	// 设置认证信息
	mac := credentials.NewCredentials(config.Config.Qiniu.AK, config.Config.Qiniu.SK)
	logger.Info(model.LogMsg{Text: "七牛云认证信息设置完成"})

	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	logger.Info(model.LogMsg{Text: "七牛云上传管理器创建完成"})

	reader := bytes.NewReader(fileData)
	err := uploadManager.UploadReader(context.Background(), reader, &uploader.ObjectOptions{
		BucketName: config.Config.Qiniu.Bucket,
		FileName:   filePath,
		ObjectName: &filePath,
	}, nil)

	if err != nil {
		logger.Error(model.LogMsg{Text: "七牛云上传失败", Data: map[string]interface{}{"err": err.Error()}})
		return "", fmt.Errorf("failed to upload file to Qiniu: %v", err)
	}
	logger.Info(model.LogMsg{Text: "七牛云上传成功"})

	return filePath, nil
}
