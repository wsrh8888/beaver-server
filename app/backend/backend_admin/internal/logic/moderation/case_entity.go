package moderation

import (
	"fmt"
	"time"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/types"
)

func formatCaseInfo(c backend_models.AdminModerationCase) types.ModerationCaseInfo {
	handleTime := ""
	if c.HandleTime != nil {
		handleTime = c.HandleTime.Format("2006-01-02 15:04:05")
	}
	return types.ModerationCaseInfo{
		ID:           uint64(c.Id),
		CaseNo:       c.CaseNo,
		Source:       c.Source,
		SourceID:     c.SourceID,
		TargetType:   c.TargetType,
		TargetID:     c.TargetID,
		Title:        c.Title,
		Description:  c.Description,
		Priority:     c.Priority,
		Status:       c.Status,
		HandlerID:    c.HandlerID,
		HandleRemark: c.HandleRemark,
		HandleTime:   handleTime,
		CreatedAt:    c.CreatedAt.String(),
	}
}

func newCaseNo() string {
	return fmt.Sprintf("CASE-%s%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)
}
