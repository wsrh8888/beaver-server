package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type AuditDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuditDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditDeveloperLogic {
	return &AuditDeveloperLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AuditDeveloperLogic) AuditDeveloper(in *open_rpc.AuditDeveloperReq) (*open_rpc.AuditDeveloperRes, error) {
	if in.Status != 1 && in.Status != 2 {
		return nil, status.Error(codes.InvalidArgument, "无效的审核状态")
	}

	var dev open_models.OpenDeveloper
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&dev).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "申请记录不存在")
		}
		return nil, err
	}
	if dev.Status != 0 {
		return nil, status.Error(codes.FailedPrecondition, "该申请已审核")
	}

	now := time.Now()
	if err := l.svcCtx.DB.Model(&dev).Updates(map[string]interface{}{
		"status":       int(in.Status),
		"audit_by":     in.AuditBy,
		"audit_time":   now.UnixMilli(),
		"audit_remark": in.AuditRemark,
		"updated_at":   now,
	}).Error; err != nil {
		l.Errorf("审核开发者失败: %v", err)
		return nil, status.Error(codes.Internal, "审核失败")
	}

	return &open_rpc.AuditDeveloperRes{}, nil
}
