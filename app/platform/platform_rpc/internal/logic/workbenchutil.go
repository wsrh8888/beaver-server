package logic

import (
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toWorkbenchAppItem(app platform_models.WorkbenchApp) *platform_rpc.WorkbenchAppItem {
	return &platform_rpc.WorkbenchAppItem{
		Id:             uint64(app.Id),
		WorkbenchAppId: app.WorkbenchAppID,
		Name:           app.Name,
		Description:    app.Description,
		Icon:           app.Icon,
		EntryUrl:       app.EntryURL,
		Category:       app.Category,
		Sort:           int32(app.Sort),
		Status:         int32(app.Status),
		CreatedBy:      app.CreatedBy,
		LastModifiedBy: app.LastModifiedBy,
		Remark:         app.Remark,
		CreatedAt:      time.Time(app.CreatedAt).Format(time.RFC3339),
		UpdatedAt:      time.Time(app.UpdatedAt).Format(time.RFC3339),
	}
}

func toWorkbenchAppPublicItem(app platform_models.WorkbenchApp) *platform_rpc.WorkbenchAppPublicItem {
	return &platform_rpc.WorkbenchAppPublicItem{
		WorkbenchAppId: app.WorkbenchAppID,
		Name:           app.Name,
		Description:    app.Description,
		Icon:           app.Icon,
		EntryUrl:       app.EntryURL,
		Category:       app.Category,
		Sort:           int32(app.Sort),
	}
}
