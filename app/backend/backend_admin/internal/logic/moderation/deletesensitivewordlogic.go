package moderation

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteSensitiveWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteSensitiveWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSensitiveWordLogic {
	return &DeleteSensitiveWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteSensitiveWordLogic) DeleteSensitiveWord(req *types.DeleteSensitiveWordReq) (resp *types.DeleteSensitiveWordRes, err error) {
	var row backend_models.AdminSensitiveWord
	if err = l.svcCtx.DB.First(&row, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("敏感词不存在")
		}
		return nil, err
	}

	if err = l.svcCtx.DB.Delete(&row).Error; err != nil {
		l.Errorf("删除敏感词失败: %v", err)
		return nil, err
	}

	l.svcCtx.RecordOperation(req.UserID, "delete_sensitive_word", "sensitive_word", row.Word, 0, "删除敏感词", "success", "")
	return &types.DeleteSensitiveWordRes{}, nil
}
