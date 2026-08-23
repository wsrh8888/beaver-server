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

package file

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/file/file_rpc/types/file_rpc"
	"beaver/utils/md5"
)

// FileTypeMapper maps file extensions to file types.
var FileTypeMapper = map[string]string{
	"jpg":  "image",
	"jpeg": "image",
	"png":  "image",
	"gif":  "image",
	"bmp":  "image",
	"webp": "image",
	"mp4":  "video",
	"avi":  "video",
	"mkv":  "video",
	"mov":  "video",
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
	"txt":  "document",
	"apk":  "apk",
	"exe":  "exe",
}

// FileUploadRequest 文件上传请求结构
type FileUploadRequest struct {
	File         multipart.File
	FileHeader   *multipart.FileHeader
	OriginalName string
	ByteData     []byte
	FileMd5      string
	FileType     string
	Suffix       string
	Size         int64
}

// ValidateAndProcessFile 验证并处理文件上传
func ValidateAndProcessFile(file multipart.File, fileHeader *multipart.FileHeader, svcCtx *svc.ServiceContext) (*FileUploadRequest, error) {
	// 文件后缀白名单验证
	originalName := fileHeader.Filename
	nameList := strings.Split(originalName, ".")
	if len(nameList) < 2 {
		return nil, errors.New("文件格式不正确")
	}
	suffix := strings.ToLower(nameList[len(nameList)-1])
	if !inList(svcCtx.Config.File.WhiteList, suffix) {
		return nil, errors.New("文件类型不在白名单中")
	}

	// 确定文件类型
	fileType := GetFileType(suffix)
	if fileType == "unknown" {
		return nil, errors.New("未知文件类型")
	}

	// 检查文件大小
	maxSize, ok := svcCtx.Config.File.MaxSize[fileType]
	if !ok {
		return nil, errors.New("配置中未找到该文件类型的最大大小")
	}
	fileSizeMB := float64(fileHeader.Size) / (1024 * 1024)
	if fileSizeMB > maxSize {
		return nil, fmt.Errorf("文件大小超过最大限制: %.2fMB", maxSize)
	}

	// 读取文件内容
	byteData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 计算文件MD5
	fileMd5 := md5.MD5(byteData)

	return &FileUploadRequest{
		File:         file,
		FileHeader:   fileHeader,
		OriginalName: originalName,
		ByteData:     byteData,
		FileMd5:      fileMd5,
		FileType:     fileType,
		Suffix:       suffix,
		Size:         fileHeader.Size,
	}, nil
}

// GetFileType 根据文件后缀获取文件类型
func GetFileType(suffix string) string {
	if fileType, ok := FileTypeMapper[suffix]; ok {
		return fileType
	}
	return "unknown"
}

// BuildLocalFileURL 生成本地文件的完整访问 URL（与 file_api 保持一致）
func BuildLocalFileURL(domain, fileKey string) string {
	if domain == "" {
		return fmt.Sprintf("/api/file/preview/%s", fileKey)
	}
	return fmt.Sprintf("%s/api/file/preview/%s", strings.TrimRight(domain, "/"), fileKey)
}

// BuildQiniuFileURL 生成七牛云文件的完整访问 URL
func BuildQiniuFileURL(domain, path string) string {
	if domain == "" || domain == "your_qiniu_domain" || path == "" {
		return ""
	}
	return fmt.Sprintf("https://%s/%s", strings.TrimRight(domain, "/"), strings.TrimLeft(path, "/"))
}

// FindFileByMd5 通过 FileRpc 按 MD5 查询文件；found=false 表示未命中（非错误）
func FindFileByMd5(ctx context.Context, fileMd5 string, svcCtx *svc.ServiceContext) (*file_rpc.FileItem, bool, error) {
	res, err := svcCtx.FileRpc.GetFileByMd5(ctx, &file_rpc.GetFileByMd5Req{Md5: fileMd5})
	if err != nil {
		return nil, false, err
	}
	if res.File == nil {
		return nil, false, nil
	}
	return res.File, true, nil
}

// inList 检查元素是否在列表中
func inList(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
