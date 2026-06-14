package logic

import (
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/types/open_rpc"
)

func toDeveloperItem(dev open_models.OpenDeveloper) *open_rpc.DeveloperItem {
	return &open_rpc.DeveloperItem{
		Id:          uint64(dev.Id),
		UserId:      dev.UserID,
		RealName:    dev.RealName,
		CompanyName: dev.CompanyName,
		Phone:       dev.Phone,
		Email:       dev.Email,
		Description: dev.Description,
		Status:      int32(dev.Status),
		AuditBy:     dev.AuditBy,
		AuditTime:   dev.AuditTime,
		AuditRemark: dev.AuditRemark,
		CreatedAt:   time.Time(dev.CreatedAt).UnixMilli(),
	}
}
