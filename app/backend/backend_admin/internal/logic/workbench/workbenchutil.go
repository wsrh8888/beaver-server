package workbench

import (
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/platform/platform_rpc/types/platform_rpc"
)

func toProtoEntryConfig(cfg *types.WorkbenchEntryConfig) *platform_rpc.WorkbenchEntryConfig {
	if cfg == nil {
		return nil
	}
	return &platform_rpc.WorkbenchEntryConfig{
		Type:   int32(cfg.Type),
		Pc:     cfg.PC,
		Mobile: cfg.Mobile,
	}
}

func fromProtoEntryConfig(cfg *platform_rpc.WorkbenchEntryConfig) types.WorkbenchEntryConfig {
	if cfg == nil {
		return types.WorkbenchEntryConfig{}
	}
	return types.WorkbenchEntryConfig{
		Type:   int(cfg.Type),
		PC:     cfg.Pc,
		Mobile: cfg.Mobile,
	}
}

func toAdminAppItem(item *platform_rpc.WorkbenchAppItem) types.GetWorkbenchAppListItem {
	return types.GetWorkbenchAppListItem{
		WorkbenchAppID: item.WorkbenchAppId,
		Name:           item.Name,
		Description:    item.Description,
		Icon:           item.Icon,
		AppType:        int(item.AppType),
		ClientScope:    int(item.ClientScope),
		EntryConfig:    fromProtoEntryConfig(item.EntryConfig),
		OpenMode:       int(item.OpenMode),
		Category:       int(item.Category),
		Sort:           int(item.Sort),
		Status:         int(item.Status),
		Remark:         item.Remark,
		CreatedBy:      item.CreatedBy,
		LastModifiedBy: item.LastModifiedBy,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}
