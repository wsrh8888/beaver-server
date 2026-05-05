package open

import (
	"context"
	"errors"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuditDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditDeveloperLogic {
	return &AuditDeveloperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuditDeveloperLogic) AuditDeveloper(req *types.AuditDeveloperReq) (resp *types.AuditDeveloperRes, err error) {
	// 1. 获取审核人 ID（从 request header，由网关注入）
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	// 2. 查询申请记录
	var developer open_models.OpenDeveloper
	err = l.svcCtx.DB.Where("id = ?", req.ID).First(&developer).Error
	if err != nil {
		return nil, errors.New("申请记录不存在")
	}

	// 3. 检查状态
	if developer.Status != 0 {
		return nil, errors.New("该申请已审核")
	}

	// 4. 验证审核状态
	if req.Status != 1 && req.Status != 2 {
		return nil, errors.New("无效的审核状态")
	}

	// 5. 更新审核信息
	now := time.Now()
	updates := map[string]interface{}{
		"status":       req.Status,
		"audit_by":     req.UserID,
		"audit_time":   now.UnixMilli(),
		"audit_remark": req.AuditRemark,
		"updated_at":   now,
	}

	err = l.svcCtx.DB.Model(&developer).Updates(updates).Error
	if err != nil {
		return nil, errors.New("审核失败")
	}

	logx.Infof("开发者审核成功: id=%s, status=%d, auditor=%s", req.ID, req.Status, req.UserID)

	return &types.AuditDeveloperRes{}, nil
}
