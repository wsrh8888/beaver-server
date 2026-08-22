package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWorkbenchAppLogic {
	return &GetWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetWorkbenchAppLogic) GetWorkbenchApp(in *platform_rpc.GetWorkbenchAppReq) (*platform_rpc.GetWorkbenchAppRes, error) {
	if strings.TrimSpace(in.WorkbenchAppId) == "" {
		return nil, status.Error(codes.InvalidArgument, "应用 ID 不能为空")
	}

	var app platform_models.WorkbenchApp
	if err := l.svcCtx.DB.Where("workbench_app_id = ?", in.WorkbenchAppId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "工作台应用不存在")
		}
		return nil, err
	}

	return &platform_rpc.GetWorkbenchAppRes{App: toWorkbenchAppItem(app)}, nil
}
