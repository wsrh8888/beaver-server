package logic

import (
	"context"
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeleteWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteWorkbenchAppLogic {
	return &DeleteWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *DeleteWorkbenchAppLogic) DeleteWorkbenchApp(in *platform_rpc.DeleteWorkbenchAppReq) (*platform_rpc.DeleteWorkbenchAppRes, error) {
	if strings.TrimSpace(in.WorkbenchAppId) == "" {
		return nil, status.Error(codes.InvalidArgument, "应用 ID 不能为空")
	}

	result := l.svcCtx.DB.Where("workbench_app_id = ?", in.WorkbenchAppId).Delete(&platform_models.WorkbenchApp{})
	if result.Error != nil {
		l.Errorf("删除工作台应用失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "删除工作台应用失败")
	}
	if result.RowsAffected == 0 {
		return nil, status.Error(codes.NotFound, "工作台应用不存在")
	}

	return &platform_rpc.DeleteWorkbenchAppRes{}, nil
}
