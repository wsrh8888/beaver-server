package document_models

import "beaver/common/models"

func (CloudDocumentSpaceMember) TableName() string {
	return "cloud_document_space_members"
}

// CloudDocumentSpaceMember 空间成员（个人空间可不建记录，仅 owner 即可）
type CloudDocumentSpaceMember struct {
	models.Model
	SpaceID   string `gorm:"column:space_id;size:64;not null;uniqueIndex:idx_space_member,priority:1" json:"spaceId"` // 关联 cloud_document_spaces.space_id
	UserID    string `gorm:"column:user_id;size:64;not null;uniqueIndex:idx_space_member,priority:2" json:"userId"` // 成员 user_id
	Role      int8   `gorm:"column:role;not null;default:1" json:"role"`                                              // 空间角色：1普通成员 2空间管理员
	InvitedBy string `gorm:"column:invited_by;size:64" json:"invitedBy,omitempty"`                                    // 邀请人 user_id
}
