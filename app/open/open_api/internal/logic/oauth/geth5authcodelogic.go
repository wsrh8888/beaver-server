package oauth

import (
	"context"
	"fmt"
	"time"

	"beaver-server/app/open/open_models"
	"beaver-server/common/models"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetH5AuthCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetH5AuthCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetH5AuthCodeLogic {
	return &GetH5AuthCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetH5AuthCode H5 免登获取 authCode
func (l *GetH5AuthCodeLogic) GetH5AuthCode(req *types.GetH5AuthCodeReq) (*types.GetH5AuthCodeRes, error) {
	// 1. 验证应用是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		return nil, fmt.Errorf("应用不存在")
	}

	// 2. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用已禁用")
	}

	// 3. 从上下文获取当前用户ID（通过 AuthMiddleware 注入）
	userID, ok := l.ctx.Value("userId").(string)
	if !ok || userID == "" {
		return nil, fmt.Errorf("未登录")
	}

	// 4. 生成临时 authCode（5分钟有效）
	authCode := models.GenerateUUID()
	expireIn := int64(300) // 5分钟
	expiresAt := time.Now().Add(time.Duration(expireIn) * time.Second)

	// 5. 存储 authCode
	h5AuthCode := open_models.OpenH5AuthCode{
		Code:      authCode,
		AppID:     req.AppID,
		UserID:    userID,
		ExpiresAt: expiresAt.Unix(),
		CreatedAt: time.Now().Unix(),
	}

	if err := l.svcCtx.DB.Create(&h5AuthCode).Error; err != nil {
		logx.Errorf("存储 H5 authCode 失败: %v", err)
		return nil, fmt.Errorf("生成授权码失败")
	}

	return &types.GetH5AuthCodeRes{
		AuthCode: authCode,
		ExpireIn: expireIn,
	}, nil
}
