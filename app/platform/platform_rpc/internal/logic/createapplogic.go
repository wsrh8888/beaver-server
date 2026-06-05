package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/platform/platform_models"
	"beaver/app/platform/platform_rpc/internal/svc"
	"beaver/app/platform/platform_rpc/types/platform_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type CreateAppLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateAppLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateAppLogic {
	return &CreateAppLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *CreateAppLogic) CreateApp(in *platform_rpc.CreateAppReq) (*platform_rpc.CreateAppRes, error) {
	var existing platform_models.UpdateApp
	if err := l.svcCtx.DB.Where("name = ?", in.Name).First(&existing).Error; err == nil {
		return nil, status.Error(codes.AlreadyExists, "应用名称已存在")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	app := platform_models.UpdateApp{
		Name:        in.Name,
		AppID:       strings.ReplaceAll(uuid.New().String(), "-", ""),
		Description: in.Description,
		IsActive:    true,
	}
	if err := l.svcCtx.DB.Create(&app).Error; err != nil {
		l.Errorf("创建应用失败: %v", err)
		return nil, status.Error(codes.Internal, "创建应用失败")
	}

	return &platform_rpc.CreateAppRes{Id: uint64(app.Id), AppId: app.AppID}, nil
}
