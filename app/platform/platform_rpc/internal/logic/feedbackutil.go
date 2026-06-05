package logic

import (
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toFeedbackItem(f platform_models.FeedbackModel) *platform_rpc.FeedbackItem {
	handleTime := ""
	if f.HandleTime != nil {
		handleTime = f.HandleTime.Format(time.RFC3339)
	}
	return &platform_rpc.FeedbackItem{
		Id:           uint64(f.Id),
		UserId:       f.UserID,
		Content:      f.Content,
		Type:         int32(f.Type),
		Status:       int32(f.Status),
		FileNames:    []string(f.FileNames),
		HandlerId:    f.HandlerID,
		HandleTime:   handleTime,
		HandleResult: f.HandleResult,
		CreatedAt:    time.Time(f.CreatedAt).Format(time.RFC3339),
		UpdatedAt:    time.Time(f.UpdatedAt).Format(time.RFC3339),
	}
}
