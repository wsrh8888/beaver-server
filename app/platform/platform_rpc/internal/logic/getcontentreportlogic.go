package logic

import (
	"context"
	"errors"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetContentReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetContentReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetContentReportLogic {
	return &GetContentReportLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetContentReportLogic) GetContentReport(in *platform_rpc.GetContentReportReq) (*platform_rpc.GetContentReportRes, error) {
	var report platform_models.ContentReportModel
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&report).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "举报不存在")
		}
		return nil, err
	}
	return &platform_rpc.GetContentReportRes{Report: toContentReportItem(report)}, nil
}
