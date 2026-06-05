package moderation

import (
	"context"
	"errors"
	"strings"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateSensitiveWordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSensitiveWordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSensitiveWordLogic {
	return &UpdateSensitiveWordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSensitiveWordLogic) UpdateSensitiveWord(req *types.UpdateSensitiveWordReq) (resp *types.UpdateSensitiveWordRes, err error) {
	var row backend_models.AdminSensitiveWord
	if err = l.svcCtx.DB.First(&row, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("敏感词不存在")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if w := strings.TrimSpace(req.Word); w != "" {
		updates["word"] = w
	}
	if req.Category != "" {
		updates["category"] = strings.TrimSpace(req.Category)
	}
	if req.Level > 0 {
		updates["level"] = req.Level
	}
	if req.Remark != "" {
		updates["remark"] = strings.TrimSpace(req.Remark)
	}
	updates["is_active"] = req.IsActive

	if len(updates) == 0 {
		return &types.UpdateSensitiveWordRes{}, nil
	}

	if err = l.svcCtx.DB.Model(&row).Updates(updates).Error; err != nil {
		l.Errorf("更新敏感词失败: %v", err)
		return nil, errors.New("更新失败")
	}

	l.svcCtx.RecordOperation(req.UserID, "update_sensitive_word", "sensitive_word", row.Word, 0, "更新敏感词", "success", "")
	return &types.UpdateSensitiveWordRes{}, nil
}
