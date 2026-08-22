package logic

import (
	"strings"
	"time"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toProtoEntryConfig(cfg platform_models.WorkbenchEntryConfig) *platform_rpc.WorkbenchEntryConfig {
	return &platform_rpc.WorkbenchEntryConfig{
		Type:   int32(cfg.Type),
		Pc:     cfg.PC,
		Mobile: cfg.Mobile,
	}
}

func fromProtoEntryConfig(cfg *platform_rpc.WorkbenchEntryConfig) platform_models.WorkbenchEntryConfig {
	if cfg == nil {
		return platform_models.WorkbenchEntryConfig{}
	}
	return platform_models.WorkbenchEntryConfig{
		Type:   int8(cfg.Type),
		PC:     strings.TrimSpace(cfg.Pc),
		Mobile: strings.TrimSpace(cfg.Mobile),
	}
}

func toWorkbenchAppItem(app platform_models.WorkbenchApp) *platform_rpc.WorkbenchAppItem {
	return &platform_rpc.WorkbenchAppItem{
		Id:             uint64(app.Id),
		WorkbenchAppId: app.WorkbenchAppID,
		Name:           app.Name,
		Description:    app.Description,
		Icon:           app.Icon,
		AppType:        int32(app.AppType),
		ClientScope:    int32(app.ClientScope),
		EntryConfig:    toProtoEntryConfig(app.EntryConfig),
		OpenMode:       int32(app.OpenMode),
		Category:       int32(app.Category),
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
		AppType:        int32(app.AppType),
		ClientScope:    int32(app.ClientScope),
		EntryConfig:    toProtoEntryConfig(app.EntryConfig),
		OpenMode:       int32(app.OpenMode),
		Category:       int32(app.Category),
		Sort:           int32(app.Sort),
	}
}

func validateEntryConfig(appType int8, cfg platform_models.WorkbenchEntryConfig) string {
	// type: 0 路由，1 URL；appType: 0 内部，1 第三方
	if cfg.Type != 0 && cfg.Type != 1 {
		return "入口类型不合法"
	}
	if appType == 0 && cfg.Type != 0 {
		return "内部应用入口类型须为路由"
	}
	if appType == 1 && cfg.Type != 1 {
		return "第三方应用入口类型须为 URL"
	}
	if strings.TrimSpace(cfg.PC) == "" && strings.TrimSpace(cfg.Mobile) == "" {
		return "PC / 移动端入口至少填写一个"
	}
	if cfg.Type == 1 {
		for _, u := range []string{cfg.PC, cfg.Mobile} {
			u = strings.TrimSpace(u)
			if u == "" {
				continue
			}
			if !strings.HasPrefix(strings.ToLower(u), "http://") && !strings.HasPrefix(strings.ToLower(u), "https://") {
				return "H5 入口须以 http:// 或 https:// 开头"
			}
		}
	}
	return ""
}
