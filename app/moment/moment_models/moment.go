package moment_models

import (
	"beaver/app/user/user_models"
	"beaver/common/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// FileInfo 结构体，用于存储文件的信息
type FileInfo struct {
	URL  string `json:"url"`  // 文件URL
	Type string `json:"type"` // 文件类型
}

type Files []FileInfo

// Value converts the Files slice to a JSON-encoded string for database storage
func (f *Files) Value() (driver.Value, error) {
	return json.Marshal(f)
}

// Scan converts a JSON-encoded string from the database to a Files slice
func (f *Files) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, f)
}

/**
 * @description: 动态表
 */
type MomentModel struct {
	models.Model
	UserID          string                `gorm:"size:64;not null" json:"userId"`             // 用户Id
	Content         string                `json:"content"`                                    // 动态内容
	Files           *Files                `gorm:"type:longtext" json:"files" `                // 文件信息（JSON数组），包括文件URL和类型
	CommentsModel   []MomentCommentModel  `gorm:"foreignkey:MomentID;references:ID" json:"-"` // 评论列表
	LikesModel      []MomentLikeModel     `gorm:"foreignkey:MomentID;references:ID" json:"-"` // 点赞列表
	MomentUserModel user_models.UserModel `gorm:"foreignKey:UserID;references:UUID" json:"-"`
	IsDeleted       bool                  `gorm:"not null;default:false" json:"isDeleted"` // 标记用户是否删除会话

}
