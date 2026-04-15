package open

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeveloperListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDeveloperListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperListLogic {
	return &GetDeveloperListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDeveloperListLogic) GetDeveloperList(req *types.GetDeveloperListReq) (resp *types.GetDeveloperListRes, err error) {
	// 1. 构建查询
	query := l.svcCtx.DB.Model(&backend_models.OpenDeveloper{})
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 2. 获取总数
	var total int64
	query.Count(&total)

	// 3. 分页查询
	var developers []backend_models.OpenDeveloper
	offset := (req.Page - 1) * req.PageSize
	query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&developers)

	// 4. 转换为响应格式
	list := make([]types.DeveloperInfo, 0, len(developers))
	for _, dev := range developers {
		list = append(list, types.DeveloperInfo{
			ID:          dev.ID,
			UserID:      dev.UserID,
			RealName:    dev.RealName,
			CompanyName: dev.CompanyName,
			Phone:       dev.Phone,
			Email:       dev.Email,
			Description: dev.Description,
			Status:      dev.Status,
			AuditBy:     dev.AuditBy,
			AuditTime:   dev.AuditTime,
			AuditRemark: dev.AuditRemark,
			CreatedAt:   dev.CreatedAt,
		})
	}

	return &types.GetDeveloperListRes{
		Total: total,
		List:  list,
	}, nil
}
