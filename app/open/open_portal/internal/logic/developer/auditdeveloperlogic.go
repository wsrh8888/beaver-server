package developer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"

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
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, errors.New("未登录")
	}

	// 验证状态
	if req.Status != 1 && req.Status != 2 {
		return nil, fmt.Errorf("无效的状态值")
	}

	// 更新开发者状态
	updates := map[string]interface{}{
		"status":       req.Status,
		"audit_time":   time.Now().Unix(),
		"audit_remark": req.AuditRemark,
		// TODO: 从上下文获取当前管理员ID
		// "audit_by":   currentAdminID,
	}

	err = l.svcCtx.DB.Model(&open_models.OpenDeveloper{}).
		Where("id = ?", req.ID).
		Updates(updates).Error

	if err != nil {
		return nil, err
	}

	l.Infof("开发者申请 %d 已审核，状态: %d", req.ID, req.Status)

	return &types.AuditDeveloperRes{}, nil
}