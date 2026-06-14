package logic

import (
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toContentReportItem(r platform_models.ContentReportModel) *platform_rpc.ContentReportItem {
	return &platform_rpc.ContentReportItem{
		Id:             uint64(r.Id),
		ReporterUserId: r.ReporterUserID,
		TargetType:     int32(r.TargetType),
		TargetId:       r.TargetID,
		ReasonType:     int32(r.ReasonType),
		Content:        r.Content,
		FileNames:      []string(r.FileNames),
		Status:         int32(r.Status),
		CaseId:         r.CaseID,
		HandlerId:      r.HandlerID,
		HandleRemark:   r.HandleRemark,
		CreatedAt:      time.Time(r.CreatedAt).Format(time.RFC3339),
	}
}
