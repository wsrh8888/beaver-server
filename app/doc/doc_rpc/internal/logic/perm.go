package logic

import (
	doc_models "beaver/app/doc/doc_models"

	"gorm.io/gorm"
)

// principalRolesFor 查询某用户在给定资源集上的显式授权角色（user 主体 + 其 group/dept 由上层扩展，本期仅 user）。
// 返回 resource_id -> role。
func explicitRole(db *gorm.DB, resourceType, resourceID, userID string) (string, error) {
	if trim(userID) == "" {
		return "", nil
	}
	var row doc_models.DocPermission
	err := db.Where("resource_type = ? AND resource_id = ? AND principal_type = ? AND principal_id = ?",
		resourceType, resourceID, "user", userID).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return row.Role, nil
}

// resolveSpaceRole 计算用户对空间的有效角色：显式授权 与 归属人（owner_user_id）取高者。
func resolveSpaceRole(db *gorm.DB, space *doc_models.DocSpace, userID string) (string, error) {
	if space == nil {
		return "", nil
	}
	role := ""
	if trim(userID) != "" && trim(space.OwnerUserID) == trim(userID) {
		role = roleOwner
	}
	exp, err := explicitRole(db, "space", space.SpaceID, userID)
	if err != nil {
		return "", err
	}
	return maxRole(role, exp), nil
}

// resolveDocRole 计算用户对文档的有效角色：
// 文档归属人 / 空间归属人 / 文档显式授权 / 空间显式授权 四者取高者。
// 归属人继承是必须的：创建者默认不写 doc_permissions 行，若不继承会出现「建完自己读不了」。
func resolveDocRole(db *gorm.DB, spaceID, docID, userID string) (string, error) {
	if trim(userID) == "" {
		return "", nil
	}
	role := ""

	var doc doc_models.DocDocument
	if err := db.Select("owner_user_id").Where("doc_id = ?", docID).First(&doc).Error; err == nil &&
		trim(doc.OwnerUserID) == trim(userID) {
		role = roleOwner
	}
	var space doc_models.DocSpace
	if err := db.Select("owner_user_id").Where("space_id = ?", spaceID).First(&space).Error; err == nil &&
		trim(space.OwnerUserID) == trim(userID) {
		role = maxRole(role, roleOwner)
	}

	docRole, err := explicitRole(db, "doc", docID, userID)
	if err != nil {
		return "", err
	}
	spaceRole, err := explicitRole(db, "space", spaceID, userID)
	if err != nil {
		return "", err
	}
	return maxRole(maxRole(role, docRole), spaceRole), nil
}

// loadDoc 取未删文档行。
func loadDoc(db *gorm.DB, docID string) (*doc_models.DocDocument, error) {
	var row doc_models.DocDocument
	if err := db.Where("doc_id = ? AND deleted_at = ?", docID, "1970-01-01 00:00:00").First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// loadSpace 取未删空间行。
func loadSpace(db *gorm.DB, spaceID string) (*doc_models.DocSpace, error) {
	var row doc_models.DocSpace
	if err := db.Where("space_id = ? AND deleted_at = ?", spaceID, "1970-01-01 00:00:00").First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// loadNode 取未删节点行。
func loadNode(db *gorm.DB, nodeID string) (*doc_models.DocNode, error) {
	var row doc_models.DocNode
	if err := db.Where("node_id = ? AND deleted_at = ?", nodeID, "1970-01-01 00:00:00").First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}
