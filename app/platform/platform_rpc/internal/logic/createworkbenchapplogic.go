package logic

import (
	"context"
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateWorkbenchAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateWorkbenchAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateWorkbenchAppLogic {
	return &CreateWorkbenchAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateWorkbenchAppLogic) CreateWorkbenchApp(in *platform_rpc.CreateWorkbenchAppReq) (*platform_rpc.CreateWorkbenchAppRes, error) {
	name := strings.TrimSpace(in.Name)
	entryURL := strings.TrimSpace(in.EntryUrl)
	if name == "" {
		return nil, status.Error(codes.InvalidArgument, "应用名称不能为空")
	}
	if entryURL == "" {
		return nil, status.Error(codes.InvalidArgument, "入口 URL 不能为空")
	}

	statusVal := int8(in.Status)
	if statusVal != 0 && statusVal != 1 {
		return nil, status.Error(codes.InvalidArgument, "状态不合法")
	}

	app := platform_models.WorkbenchApp{
		WorkbenchAppID: strings.ReplaceAll(uuid.New().String(), "-", ""),
		Name:           name,
		Description:    strings.TrimSpace(in.Description),
		Icon:           strings.TrimSpace(in.Icon),
		EntryURL:       entryURL,
		Category:       strings.TrimSpace(in.Category),
		Sort:           int(in.Sort),
		Status:         statusVal,
		Remark:         strings.TrimSpace(in.Remark),
		CreatedBy:      in.OperatorId,
		LastModifiedBy: in.OperatorId,
	}
	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		l.Errorf("创建工作台应用失败: %v", err)
		return nil, status.Error(codes.Internal, "创建工作台应用失败")
	}

	return &platform_rpc.CreateWorkbenchAppRes{WorkbenchAppId: app.WorkbenchAppID}, nil
}
