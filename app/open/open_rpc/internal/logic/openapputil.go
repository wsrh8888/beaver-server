package logic

import (
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"
)

func toOpenAppItem(app open_models.OpenApp) *open_rpc.OpenAppItem {
	var auditedAt int64
	if app.AuditedAt != nil {
		auditedAt = app.AuditedAt.UnixMilli()
	}
	return &open_rpc.OpenAppItem{
		AppId:         app.AppID,
		Name:          app.Name,
		Description:   app.Description,
		Icon:          app.Icon,
		OwnerUserId:   app.OwnerUserID,
		AppType:       int32(app.AppType),
		Category:      app.Category,
		Status:        int32(app.Status),
		AuditStatus:   int32(app.AuditStatus),
		AuditedBy:     app.AuditedBy,
		AuditedAt:     auditedAt,
		EnableRobot:   int32(app.EnableRobot),
		EnableOauth:   int32(app.EnableOAuth),
		EnableWebhook: int32(app.EnableWebhook),
		CreatedAt:     time.Time(app.CreatedAt).UnixMilli(),
	}
}
