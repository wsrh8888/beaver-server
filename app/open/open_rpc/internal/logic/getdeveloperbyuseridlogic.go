package logic

import (
	"context"
	"errors"

	"beaver/app/open/open_models"
	"beaver/app/open/open_rpc/internal/svc"
	"beaver/app/open/open_rpc/types/open_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetDeveloperByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeveloperByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeveloperByUserIDLogic {
	return &GetDeveloperByUserIDLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *GetDeveloperByUserIDLogic) GetDeveloperByUserID(in *open_rpc.GetDeveloperByUserIDReq) (*open_rpc.GetDeveloperByUserIDRes, error) {
	var dev open_models.OpenDeveloper
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&dev).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &open_rpc.GetDeveloperByUserIDRes{Found: false}, nil
	}
	if err != nil {
		return nil, err
	}
	return &open_rpc.GetDeveloperByUserIDRes{
		Found:     true,
		Developer: toDeveloperItem(dev),
	}, nil
}
