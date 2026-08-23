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
	filecommon "beaver/app/backend/backend_admin/internal/handler/file/common"
	logic "beaver/app/backend/backend_admin/internal/logic/file"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	utils "beaver/utils/list"
	"beaver/utils/md5"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileUploadQiniuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("开始处理文件上传请求")

		var req types.FileUploadQiniuReq
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
			logx.Error("文件名格式不正确:", originalName)
			response.Response(r, w, nil, errors.New("文件格式不正确"))
			return
		}
		suffix := strings.ToLower(nameList[len(nameList)-1])
		if !utils.InList(svcCtx.Config.File.WhiteList, suffix) {
			logx.Error("文件类型不在白名单中:", suffix)
			response.Response(r, w, nil, errors.New("文件类型不支持"))
			return
		}
		logx.Info("文件类型检查通过:", suffix)

		// 确定文件类型
		fileType := filecommon.GetFileType(suffix)
		if fileType == "unknown" {
			logx.Error("未知的文件类型:", suffix)
			response.Response(r, w, nil, errors.New("不支持的文件类型"))
			return
		}
		logx.Info("文件类型:", fileType)

		// 检查文件大小
		maxSize, ok := svcCtx.Config.File.MaxSize[fileType]
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
		existingFile, found, err := filecommon.FindFileByMd5(r.Context(), fileMd5, svcCtx)
		if err != nil {
			logx.Error("查询文件失败:", err)
			response.Response(r, w, nil, errors.New("查询文件失败"))
			return
		}
		if found {
			logx.Info("文件已存在，直接返回:", existingFile.FileKey, existingFile.OriginalName)
			resp.OriginalName = existingFile.OriginalName
			resp.FileURL = filecommon.BuildQiniuFileURL(svcCtx.Config.Qiniu.Domain, existingFile.Path)
			response.Response(r, w, resp, nil)
			return
		}
		logx.Info("文件不存在，继续上传流程")

		if fileInfoStr == "" {
			logx.Error("fileInfo不能为空")
			response.Response(r, w, nil, errors.New("fileInfo不能为空"))
			return
		}
		if !json.Valid([]byte(fileInfoStr)) {
			logx.Error("fileInfo格式不正确")
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
			logx.Error("保存文件信息失败:", err)
			response.Response(r, w, nil, errors.New("保存文件信息失败"))
			return
		}
		logx.Info("数据库记录创建成功:", saveResp.FileKey)

		resp.OriginalName = saveReq.OriginalName
		resp.FileURL = filecommon.BuildQiniuFileURL(svcCtx.Config.Qiniu.Domain, qiniuURL)

		logx.Info("文件上传完成, url:", resp.FileURL)
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
