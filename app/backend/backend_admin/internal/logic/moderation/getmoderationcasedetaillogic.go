package moderation

import (
	"context"
	"errors"

	"beaver/app/backend/backend_models"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetModerationCaseDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetModerationCaseDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetModerationCaseDetailLogic {
	return &GetModerationCaseDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetModerationCaseDetailLogic) GetModerationCaseDetail(req *types.GetModerationCaseDetailReq) (resp *types.GetModerationCaseDetailRes, err error) {
	if req.CaseID == 0 {
		return nil, errors.New("工单ID不能为空")
	}

	var c backend_models.AdminModerationCase
	if err = l.svcCtx.DB.Where("id = ?", req.CaseID).First(&c).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("工单不存在")
		}
		l.Errorf("查询工单详情失败: %v", err)
		return nil, err
	}
	return &types.GetModerationCaseDetailRes{Case: formatCaseInfo(c)}, nil
}
