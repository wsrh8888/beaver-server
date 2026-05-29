package logic

import (
	"errors"
	"fmt"

	"beaver/app/document/document_models"
	"beaver/common/models/ctype"
	"gorm.io/gorm"
)

func resolveDocumentPerm(db *gorm.DB, doc *document_models.CloudDocument, userID string) (int, error) {
	if doc.OwnerID == userID {
		return document_models.DocPermManage, nil
	}
	var row document_models.CloudDocumentPermission
	err := db.Where("doc_id = ? AND subject_type = ? AND subject_id = ?",
		doc.DocID, document_models.DocSubjectUser, userID).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return row.Perm, nil
}

func canRead(perm int) bool {
	return perm >= document_models.DocPermView
}

func canEdit(perm int) bool {
	return perm >= document_models.DocPermEdit
}

func canManage(perm int) bool {
	return perm >= document_models.DocPermManage
}

func loadActiveDocument(db *gorm.DB, docID string) (*document_models.CloudDocument, error) {
	var doc document_models.CloudDocument
	err := db.Where("doc_id = ? AND status = 1 AND deleted_at IS NULL", docID).First(&doc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("文档不存在")
	}
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func resolveSpaceAccess(db *gorm.DB, userID, reqSpaceID string) (string, error) {
	if reqSpaceID == "" {
		return ensureUserSpace(db, userID)
	}
	var space document_models.CloudDocumentSpace
	err := db.Where("space_id = ? AND status = 1", reqSpaceID).First(&space).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("空间不存在")
	}
	if err != nil {
		return "", err
	}
	if space.OwnerID == userID {
		return reqSpaceID, nil
	}
	var member document_models.CloudDocumentSpaceMember
	err = db.Where("space_id = ? AND user_id = ?", reqSpaceID, userID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", fmt.Errorf("无权限访问该空间")
	}
	if err != nil {
		return "", err
	}
	return reqSpaceID, nil
}

func ensureSpaceReadAccess(db *gorm.DB, userID, spaceID string) error {
	var space document_models.CloudDocumentSpace
	err := db.Where("space_id = ? AND status = 1", spaceID).First(&space).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("空间不存在")
	}
	if err != nil {
		return err
	}
	if space.OwnerID == userID {
		return nil
	}
	var member document_models.CloudDocumentSpaceMember
	err = db.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("无权限访问该空间")
	}
	if err != nil {
		return err
	}
	if space.DefaultPerm < document_models.DocPermView {
		return fmt.Errorf("无权限访问该空间")
	}
	return nil
}

func ensureSpaceWriteAccess(db *gorm.DB, userID, spaceID string) error {
	var space document_models.CloudDocumentSpace
	err := db.Where("space_id = ? AND status = 1", spaceID).First(&space).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("空间不存在")
	}
	if err != nil {
		return err
	}
	if space.OwnerID == userID {
		return nil
	}
	var member document_models.CloudDocumentSpaceMember
	err = db.Where("space_id = ? AND user_id = ?", spaceID, userID).First(&member).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("无权限访问该空间")
	}
	if err != nil {
		return err
	}
	if space.DefaultPerm < document_models.DocPermEdit {
		return fmt.Errorf("无权限在该空间下操作")
	}
	return nil
}

func checkParentReadAccess(db *gorm.DB, userID, spaceID, parentID string) error {
	if err := ensureSpaceReadAccess(db, userID, spaceID); err != nil {
		return err
	}
	if parentID == "" {
		return nil
	}
	parent, err := loadActiveDocument(db, parentID)
	if err != nil {
		return err
	}
	if parent.SpaceID != spaceID {
		return fmt.Errorf("父文件夹不属于该空间")
	}
	if parent.DocType != ctype.CloudDocTypeFolder {
		return fmt.Errorf("父节点不是文件夹")
	}
	perm, err := resolveDocumentPerm(db, parent, userID)
	if err != nil {
		return err
	}
	if !canRead(perm) {
		return fmt.Errorf("无权限访问该目录")
	}
	return nil
}

func checkParentWriteAccess(db *gorm.DB, userID, spaceID, parentID string) error {
	if err := ensureSpaceWriteAccess(db, userID, spaceID); err != nil {
		return err
	}
	if parentID == "" {
		return nil
	}
	parent, err := loadActiveDocument(db, parentID)
	if err != nil {
		return err
	}
	if parent.SpaceID != spaceID {
		return fmt.Errorf("父文件夹不属于该空间")
	}
	if parent.DocType != ctype.CloudDocTypeFolder {
		return fmt.Errorf("父节点不是文件夹")
	}
	perm, err := resolveDocumentPerm(db, parent, userID)
	if err != nil {
		return err
	}
	if !canEdit(perm) {
		return fmt.Errorf("无权限在该目录下操作")
	}
	return nil
}

func isDescendantOf(db *gorm.DB, ancestorID, nodeID string) (bool, error) {
	if ancestorID == nodeID {
		return true, nil
	}
	current := nodeID
	for current != "" {
		doc, err := loadActiveDocument(db, current)
		if err != nil {
			return false, err
		}
		if doc.ParentID == ancestorID {
			return true, nil
		}
		current = doc.ParentID
	}
	return false, nil
}
