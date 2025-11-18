package file_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FileType 文件类型
type FileType string

const (
	ImageFileType    FileType = "image"
	VideoFileType    FileType = "video"
	AudioFileType    FileType = "audio"
	DocumentFileType FileType = "document"
	ArchiveFileType  FileType = "archive"
	OtherFileType    FileType = "other"
)

// ImageFile 图片文件信息
type ImageFile struct {
	Width  int `json:"width"`  // 图片宽度
	Height int `json:"height"` // 图片高度
}

// VideoFile 视频文件信息
type VideoFile struct {
	Width    int `json:"width"`    // 视频宽度
	Height   int `json:"height"`   // 视频高度
	Duration int `json:"duration"` // 视频时长（秒）
}

// AudioFile 音频文件信息
type AudioFile struct {
	Duration int `json:"duration"` // 音频时长（秒）
}

// FileInfo 文件信息，类似ctype.Msg的结构
type FileInfo struct {
	Type      FileType   `json:"type"`                // 文件类型
	ImageFile *ImageFile `json:"imageFile,omitempty"` // 图片文件信息
	VideoFile *VideoFile `json:"videoFile,omitempty"` // 视频文件信息
	AudioFile *AudioFile `json:"audioFile,omitempty"` // 音频文件信息
}

// Value converts the FileInfo to a JSON-encoded string for database storage
func (f *FileInfo) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan converts a JSON-encoded string from the database to a FileInfo
func (f *FileInfo) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

// FileSource 文件来源
type FileSource string

const (
	QiniuSource FileSource = "qiniu" // 七牛云
	LocalSource FileSource = "local" // 本地存储
)

type FileModel struct {
	models.Model
	FileKey      string     `json:"fileKey"`                       // 文件唯一ID /api/file/{uuid}
	OriginalName string     `json:"originalName"`                  // 原始文件名（带后缀名）
	Size         int64      `json:"size"`                          // 文件大小
	Path         string     `json:"path"`                          // 文件实际存储路径
	Md5          string     `json:"md5"`                           // 文件md5
	Type         string     `json:"type"`                          // 文件类型
	Source       FileSource `json:"source" default:"qiniu"`        // 文件来源：qiniu(七牛云) 或 local(本地)
	FileInfo     *FileInfo  `gorm:"type:longtext" json:"fileInfo"` // 文件详细信息（JSON格式）
}
