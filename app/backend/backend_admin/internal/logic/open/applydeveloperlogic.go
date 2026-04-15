package open

import (
	"context"
	"errors"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/backend/backend_models"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyDeveloperLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyDeveloperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyDeveloperLogic {
	return &ApplyDeveloperLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyDeveloperLogic) ApplyDeveloper(req *types.ApplyDeveloperReq) (resp *types.ApplyDeveloperRes, err error) {
	// 1. 从 header 获取当前用户 ID（需要在 handler 里注入）
	userID := l.ctx.Value("userId")
	if userID == nil {
		return nil, errors.New("未登录")
	}

	// 2. 检查是否已经申请过
	var existing backend_models.OpenDeveloper
	err = l.svcCtx.DB.Where("user_id = ?", userID).First(&existing).Error
	if err == nil {
		return nil, errors.New("您已经提交过申请，请等待审核")
	}

	// 3. 创建申请记录
	id := uuid.New().String()
	developer := backend_models.OpenDeveloper{
		Model: backend_models.Model{
			ID:        id,
			CreatedAt: time.Now().UnixMilli(),
			UpdatedAt: time.Now().UnixMilli(),
		},
		UserID:      userID.(string),
		RealName:    req.RealName,
		CompanyName: req.CompanyName,
		Phone:       req.Phone,
		Email:       req.Email,
		Description: req.Description,
		Status:      0, // 待审核
	}

	err = l.svcCtx.DB.Create(&developer).Error
	if err != nil {
		return nil, errors.New("申请失败")
	}

	return &types.ApplyDeveloperRes{
		ID: id,
	}, nil
}
