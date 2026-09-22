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

package common

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"

	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	utils "beaver/utils/list"
	"beaver/utils/md5"
)

// minioPresignExpiry MinIO预签名URL有效期
const minioPresignExpiry = 3600 // 秒

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
	"m4a":  "audio",
	"aac":  "audio",
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
}

// FileUploadRequest 文件上传请求结构
type FileUploadRequest struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
	ByteData   []byte
	FileMd5    string
	FileType   string
	Suffix     string
	Size       int64
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
	if !utils.InList(svcCtx.Config.WhiteList, suffix) {
		return nil, errors.New("文件类型不在白名单中")
	}

	// 确定文件类型
	fileType := getFileType(suffix)
	if fileType == "unknown" {
		return nil, errors.New("未知文件类型")
	}

	// 检查文件大小
	maxSize, ok := svcCtx.Config.FileMaxSize[fileType]
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
		File:       file,
		FileHeader: fileHeader,
		ByteData:   byteData,
		FileMd5:    fileMd5,
		FileType:   fileType,
		Suffix:     suffix,
		Size:       fileHeader.Size,
	}, nil
}

// CheckFileExists 检查文件是否已存在于数据库中
func CheckFileExists(fileMd5 string, svcCtx *svc.ServiceContext) (*file_models.FileModel, error) {
	var fileModel file_models.FileModel
	err := svcCtx.DB.Take(&fileModel, "md5 = ?", fileMd5).Error
	if err != nil {
		return nil, err
	}
	return &fileModel, nil
}

// CreateFileRecord 创建文件记录
func CreateFileRecord(req *FileUploadRequest, filePath string, source file_models.FileSource, svcCtx *svc.ServiceContext) (*file_models.FileModel, error) {
	FileKeyWithSuffix := req.FileMd5 + "." + req.Suffix

	// 创建文件记录
	newFileModel := &file_models.FileModel{
		OriginalName: strings.TrimSuffix(req.FileHeader.Filename, "."+req.Suffix),
		Size:         req.Size,
		Path:         filePath,
		Md5:          req.FileMd5,
		FileKey:      FileKeyWithSuffix,
		Type:         req.FileType,
		Source:       source,
	}

	// 保存到数据库
	err := svcCtx.DB.Create(newFileModel).Error
	if err != nil {
		return nil, fmt.Errorf("保存文件记录失败: %v", err)
	}

	beaverlog.New("file_utils").Info(model.LogMsg{Text: "文件记录创建成功", Data: map[string]interface{}{"fileKey": newFileModel.FileKey, "source": source}})
	return newFileModel, nil
}

// getFileType 根据文件后缀获取文件类型
func getFileType(suffix string) string {
	if fileType, ok := FileTypeMapper[suffix]; ok {
		return fileType
	}
	return "unknown"
}

// GetFileTypeFromMimeType 从MIME类型获取文件类型（公共函数）
func GetFileTypeFromMimeType(mimeType string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "image"
	} else if strings.HasPrefix(mimeType, "video/") {
		return "video"
	} else if strings.HasPrefix(mimeType, "audio/") {
		return "audio"
	} else if strings.HasPrefix(mimeType, "application/") {
		return "document"
	}
	return "other"
}

// GetContentType 根据文件类型获取Content-Type
func GetContentType(fileType string) string {
	switch fileType {
	case "image":
		return "image/*"
	case "video":
		return "video/*"
	case "audio":
		return "audio/*"
	case "document":
		return "application/pdf"
	case "archive":
		return "application/zip"
	default:
		return "application/octet-stream"
	}
}

// ConvertFileInfoToAPI 转换FileInfo为API响应格式
func ConvertFileInfoToAPI(fileInfo *file_models.FileInfo) *types.FileInfo {
	if fileInfo == nil {
		return nil
	}

	result := &types.FileInfo{
		Type: string(fileInfo.Type),
	}

	if fileInfo.ImageFile != nil {
		result.ImageFile = &types.ImageFile{
			Width:  fileInfo.ImageFile.Width,
			Height: fileInfo.ImageFile.Height,
		}
	}

	if fileInfo.VideoFile != nil {
		result.VideoFile = &types.VideoFile{
			Width:    fileInfo.VideoFile.Width,
			Height:   fileInfo.VideoFile.Height,
			Duration: fileInfo.VideoFile.Duration,
		}
	}

	if fileInfo.AudioFile != nil {
		result.AudioFile = &types.AudioFile{
			Duration: fileInfo.AudioFile.Duration,
		}
	}

	return result
}

// ConvertAPIFileInfoToModel 转换API FileInfo为数据库模型格式
func ConvertAPIFileInfoToModel(apiFileInfo *types.FileInfo) *file_models.FileInfo {
	if apiFileInfo == nil {
		return nil
	}

	result := &file_models.FileInfo{
		Type: file_models.FileType(apiFileInfo.Type),
	}

	if apiFileInfo.ImageFile != nil {
		result.ImageFile = &file_models.ImageFile{
			Width:  apiFileInfo.ImageFile.Width,
			Height: apiFileInfo.ImageFile.Height,
		}
	}

	if apiFileInfo.VideoFile != nil {
		result.VideoFile = &file_models.VideoFile{
			Width:    apiFileInfo.VideoFile.Width,
			Height:   apiFileInfo.VideoFile.Height,
			Duration: apiFileInfo.VideoFile.Duration,
		}
	}

	if apiFileInfo.AudioFile != nil {
		result.AudioFile = &file_models.AudioFile{
			Duration: apiFileInfo.AudioFile.Duration,
		}
	}

	return result
}

// SaveFileToLocal 保存文件到本地（公共函数）
func SaveFileToLocal(filePath string, data []byte) error {
	// 确保目录存在
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	// 保存文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	beaverlog.New("file_utils").Info(model.LogMsg{Text: "文件保存成功", Data: map[string]interface{}{"filePath": filePath}})
	return nil
}

// GenerateFilePath 生成文件路径（公共函数）
// projectName: 项目名称，如果为空则不加项目目录前缀
func GenerateFilePath(uploadDir, projectName, fileType, fileMd5, suffix string) string {
	fileMd5Name := fileMd5 + "." + suffix
	if projectName != "" {
		return filepath.Join(uploadDir, projectName, fileType, fileMd5Name)
	}
	return filepath.Join(uploadDir, fileType, fileMd5Name)
}

// GenerateRelativePath 生成相对路径（不包含uploadDir，用于数据库存储）
// projectName: 项目名称，如果为空则不加项目目录前缀
func GenerateRelativePath(projectName, fileType, fileMd5, suffix string) string {
	fileMd5Name := fileMd5 + "." + suffix
	var path string
	if projectName != "" {
		path = filepath.Join(projectName, fileType, fileMd5Name)
	} else {
		path = filepath.Join(fileType, fileMd5Name)
	}
	// 使用 filepath.Join 生成路径，然后转换为正斜杠格式，确保跨平台一致性
	return filepath.ToSlash(path)
}

// ParseDuration 解析时长字符串为秒数（公共函数）
func ParseDuration(durationStr string) int {
	// 时长格式可能是 "123.456" 秒
	if durationStr == "" {
		return 0
	}

	// 简单处理，取整数部分
	if idx := strings.Index(durationStr, "."); idx != -1 {
		durationStr = durationStr[:idx]
	}

	var duration int
	if _, err := fmt.Sscanf(durationStr, "%d", &duration); err == nil {
		return duration
	}

	return 0
}

// ============================================================
// 存储后端抽象层
// ============================================================

// Storage 存储后端抽象接口，各实现只负责"把字节流存到某处"这一件事
type Storage interface {
	// Upload 上传文件，objectKey 为相对路径（如 beaver/image/md5.jpg），返回实际存储路径
	Upload(ctx context.Context, data []byte, objectKey string) (storedPath string, err error)
}

// LocalStorage 本地磁盘存储
type LocalStorage struct {
	UploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	return &LocalStorage{UploadDir: uploadDir}
}

func (s *LocalStorage) Upload(ctx context.Context, data []byte, objectKey string) (string, error) {
	filePath := filepath.Join(s.UploadDir, objectKey)
	if err := SaveFileToLocal(filePath, data); err != nil {
		return "", err
	}
	// 本地存储的 Path 字段存相对路径（objectKey），预览时再拼接 UploadDir
	return objectKey, nil
}

// QiniuStorage 七牛云对象存储
type QiniuStorage struct {
	AK     string
	SK     string
	Bucket string
}

func NewQiniuStorage(ak, sk, bucket string) *QiniuStorage {
	return &QiniuStorage{AK: ak, SK: sk, Bucket: bucket}
}

func (s *QiniuStorage) Upload(ctx context.Context, data []byte, objectKey string) (string, error) {
	mac := credentials.NewCredentials(s.AK, s.SK)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{Credentials: mac},
	})
	reader := bytes.NewReader(data)
	err := uploadManager.UploadReader(ctx, reader, &uploader.ObjectOptions{
		BucketName: s.Bucket,
		FileName:   objectKey,
		ObjectName: &objectKey,
	}, nil)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Qiniu: %v", err)
	}
	return objectKey, nil
}

// MinioStorage MinIO对象存储
type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinioStorage(client *minio.Client, bucket string) *MinioStorage {
	return &MinioStorage{Client: client, Bucket: bucket}
}

func (s *MinioStorage) Upload(ctx context.Context, data []byte, objectKey string) (string, error) {
	reader := bytes.NewReader(data)
	_, err := s.Client.PutObject(ctx, s.Bucket, objectKey, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to MinIO: %v", err)
	}
	return objectKey, nil
}

// DefaultStorage 根据yaml StorageType配置返回默认存储后端
func DefaultStorage(svcCtx *svc.ServiceContext) (Storage, file_models.FileSource, error) {
	switch svcCtx.Config.StorageType {
	case "qiniu":
		return NewQiniuStorage(svcCtx.Config.Qiniu.AK, svcCtx.Config.Qiniu.SK, svcCtx.Config.Qiniu.Bucket), file_models.QiniuSource, nil
	case "minio":
		if svcCtx.MinioClient == nil {
			return nil, "", errors.New("MinIO未配置")
		}
		return NewMinioStorage(svcCtx.MinioClient, svcCtx.Config.Minio.Bucket), file_models.MinioSource, nil
	default: // "local" 或空值
		return NewLocalStorage(svcCtx.Config.Local.UploadDir), file_models.LocalSource, nil
	}
}

// buildFileURL 按文件来源拼接对外访问URL
func buildFileURL(svcCtx *svc.ServiceContext, source file_models.FileSource, m *file_models.FileModel) string {
	switch source {
	case file_models.QiniuSource:
		domain := svcCtx.Config.Qiniu.Domain
		if domain != "" && domain != "your_qiniu_domain" {
			return fmt.Sprintf("https://%s/%s", domain, m.Path)
		}
		return ""
	case file_models.MinioSource:
		// 走预览端点，由 previewhandler 生成 MinIO 预签名URL重定向
		if svcCtx.Config.Domain != "" {
			return fmt.Sprintf("%s/api/file/preview/%s", svcCtx.Config.Domain, m.FileKey)
		}
		return fmt.Sprintf("/api/file/preview/%s", m.FileKey)
	case file_models.LocalSource:
		if svcCtx.Config.Domain != "" {
			return fmt.Sprintf("%s/api/file/preview/%s", svcCtx.Config.Domain, m.FileKey)
		}
		return fmt.Sprintf("/api/file/preview/%s", m.FileKey)
	}
	return ""
}

// projectNameFor 按来源取对应的项目名前缀
func projectNameFor(source file_models.FileSource, svcCtx *svc.ServiceContext) string {
	switch source {
	case file_models.QiniuSource:
		return svcCtx.Config.Qiniu.ProjectName
	case file_models.MinioSource:
		return svcCtx.Config.Minio.ProjectName
	default:
		return svcCtx.Config.Local.ProjectName
	}
}

// UploadFile 共享上传流程：校验 -> MD5去重 -> 存储 -> 建DB记录 -> 拼URL
// handler 负责从 *http.Request 解析 multipart.File/fileHeader/fileInfoStr，
// logic 层调用本函数完成业务编排，避免 logic 直接依赖 *http.Request
func UploadFile(ctx context.Context, svcCtx *svc.ServiceContext, file multipart.File, fileHead *multipart.FileHeader, fileInfoStr string, store Storage, source file_models.FileSource) (*types.FileRes, error) {
	// 校验并处理文件
	fileReq, err := ValidateAndProcessFile(file, fileHead, svcCtx)
	if err != nil {
		return nil, err
	}

	// MD5去重：命中则直接返回已存记录
	existingFile, err := CheckFileExists(fileReq.FileMd5, svcCtx)
	if err == nil {
		resp := &types.FileRes{
			FileKey:      existingFile.FileKey,
			OriginalName: existingFile.OriginalName,
			FileURL:      buildFileURL(svcCtx, source, existingFile),
		}
		if existingFile.FileInfo != nil {
			resp.FileInfo = ConvertFileInfoToAPI(existingFile.FileInfo)
		}
		return resp, nil
	}

	// 生成相对路径作为 objectKey
	objectKey := GenerateRelativePath(projectNameFor(source, svcCtx), fileReq.FileType, fileReq.FileMd5, fileReq.Suffix)

	// 存储到指定后端
	storedPath, err := store.Upload(ctx, fileReq.ByteData, objectKey)
	if err != nil {
		return nil, err
	}

	// 创建文件记录
	newFileModel, err := CreateFileRecord(fileReq, storedPath, source, svcCtx)
	if err != nil {
		return nil, err
	}

	// 解析FormData中的fileInfo字段
	var fileInfo *file_models.FileInfo
	if fileInfoStr != "" {
		var apiFileInfo types.FileInfo
		if json.Unmarshal([]byte(fileInfoStr), &apiFileInfo) == nil {
			fileInfo = ConvertAPIFileInfoToModel(&apiFileInfo)
		}
	}
	if fileInfo != nil {
		newFileModel.FileInfo = fileInfo
		svcCtx.DB.Save(newFileModel)
	}

	resp := &types.FileRes{
		FileKey:      newFileModel.FileKey,
		OriginalName: newFileModel.OriginalName,
		FileURL:      buildFileURL(svcCtx, source, newFileModel),
	}
	if fileInfo != nil {
		resp.FileInfo = ConvertFileInfoToAPI(fileInfo)
	}

	beaverlog.New("file_upload").Info(model.LogMsg{Text: "文件上传成功", Data: map[string]interface{}{"fileKey": newFileModel.FileKey, "source": string(source)}})
	return resp, nil
}

// PresignMinioURL 生成MinIO预签名访问URL（供 previewhandler 调用）
func PresignMinioURL(ctx context.Context, client *minio.Client, bucket, objectKey string) (string, error) {
	u, err := client.PresignedGetObject(ctx, bucket, objectKey, minioPresignExpiry*time.Second, nil)
	if err != nil {
		return "", fmt.Errorf("生成MinIO预签名URL失败: %v", err)
	}
	return u.String(), nil
}
