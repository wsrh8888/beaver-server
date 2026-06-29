package circle_models

import (
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// PostFileInfo 帖子附件信息
type PostFileInfo struct {
	FileKey string `json:"fileKey"` // 文件key
	Type    uint32 `json:"type"`    // 文件类型：2=图片 3=视频 4=文件
}

type PostFiles []PostFileInfo

func (f *PostFiles) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *PostFiles) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

// CirclePostModel 圈子帖子表
type CirclePostModel struct {
	models.Model
	PostID       string     `gorm:"column:post_id;size:64;uniqueIndex;not null" json:"postId"`  // 帖子唯一ID
	CircleID     string     `gorm:"size:64;not null;index" json:"circleId"`                     // 所属圈子ID
	UserID       string     `gorm:"size:64;not null;index" json:"userId"`                       // 发帖用户ID
	Title        string     `gorm:"size:128" json:"title"`                                      // 帖子标题（可选）
	Content      string     `gorm:"type:text;not null" json:"content"`                          // 帖子内容
	Files        *PostFiles `gorm:"type:longtext" json:"files"`                                 // 附件列表（JSON数组）
	CommentCount int64      `gorm:"not null;default:0" json:"commentCount"`                     // 评论数（冗余字段）
	LikeCount    int64      `gorm:"not null;default:0" json:"likeCount"`                        // 点赞数（冗余字段）
	IsTop        bool       `gorm:"not null;default:false" json:"isTop"`                        // 是否置顶
	IsDeleted    bool       `gorm:"not null;default:false;index" json:"isDeleted"`              // 软删除标记
}
