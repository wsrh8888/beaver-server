package open

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/open/open_models"

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
	// 1. 获取当前用户 ID（从 request header，由网关注入）
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}

	// 2. 检查是否已经申请过
	var existing open_models.OpenDeveloper
	err = l.svcCtx.DB.Where("user_id = ?", req.UserID).First(&existing).Error
	if err == nil {
		return nil, errors.New("您已经提交过申请，请等待审核")
	}

	// 3. 创建申请记录
	id := uuid.New().String()
	developer := open_models.OpenDeveloper{
		UserID:      req.UserID,
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

	logx.Infof("开发者申请成功: user_id=%s, id=%s", req.UserID, id)

	return &types.ApplyDeveloperRes{
		ID: id,
	}, nil
}
