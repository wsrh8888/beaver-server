package developer

import (
	"context"
	"time"

	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"beaver/app/open/open_models"

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
	// 构建查询
	query := l.svcCtx.DB.Model(&open_models.OpenDeveloper{})

	// 状态筛选
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	var total int64
	query.Count(&total)

	// 分页查询
	var developers []open_models.OpenDeveloper
	offset := (req.Page - 1) * req.PageSize
	err = query.Order("created_at DESC").Offset(offset).Limit(req.PageSize).Find(&developers).Error

	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	list := make([]types.DeveloperInfo, 0, len(developers))
	for _, dev := range developers {
		list = append(list, types.DeveloperInfo{
			ID:          dev.Id,
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
			CreatedAt:   time.Time(dev.CreatedAt).Unix(),
		})
	}

	return &types.GetDeveloperListRes{
		Total: total,
		List:  list,
	}, nil
}
