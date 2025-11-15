package file

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
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
	originalName := fileHeader.FileName
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

// inList 检查元素是否在列表中
func inList(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
