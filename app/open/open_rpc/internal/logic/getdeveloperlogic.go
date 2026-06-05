package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GetDeveloperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperLogic {
	return &GetDeveloperLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetDeveloperLogic) GetDeveloper(in *open_rpc.GetDeveloperReq) (*open_rpc.GetDeveloperRes, error) {
	var dev open_models.OpenDeveloper
	if err := l.svcCtx.DB.Where("id = ?", in.Id).First(&dev).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "开发者记录不存在")
		}
		return nil, err
	}
	return &open_rpc.GetDeveloperRes{Developer: toDeveloperItem(dev)}, nil
}
