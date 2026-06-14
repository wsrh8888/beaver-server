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

type DeleteBotLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteBotLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteBotLogic {
	return &DeleteBotLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteBotLogic) DeleteBot(in *open_rpc.DeleteBotReq) (*open_rpc.DeleteBotRes, error) {
	if in.Id == 0 {
		return nil, errors.New("id 不能为空")
	}

	result := l.svcCtx.DB.Model(&open_models.OpenBotModel{}).
		Where("id = ?", in.Id).
		Update("status", 0)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &open_rpc.DeleteBotRes{}, nil
}
