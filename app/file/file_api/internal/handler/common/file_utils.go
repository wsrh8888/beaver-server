package common

import (
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

	"beaver/app/file/file_api/internal/svc"
	"beaver/app/file/file_api/internal/types"
	"beaver/app/file/file_models"
	"beaver/utils"
	"beaver/utils/md5"

	"github.com/zeromicro/go-zero/core/logx"
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

	logx.Infof("文件记录创建成功: %s, 来源: %s", newFileModel.FileKey, source)
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

	logx.Infof("文件保存成功: %s", filePath)
	return nil
}

// GenerateFilePath 生成文件路径（公共函数）
func GenerateFilePath(uploadDir, fileType, fileMd5, suffix string) string {
	fileMd5Name := fileMd5 + "." + suffix
	return filepath.Join(uploadDir, fileType, fileMd5Name)
}

// GenerateRelativePath 生成相对路径（不包含uploadDir，用于数据库存储）
func GenerateRelativePath(fileType, fileMd5, suffix string) string {
	fileMd5Name := fileMd5 + "." + suffix
	// 使用 filepath.Join 生成路径，然后转换为正斜杠格式，确保跨平台一致性
	path := filepath.Join(fileType, fileMd5Name)
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
