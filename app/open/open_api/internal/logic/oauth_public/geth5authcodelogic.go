// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package oauth_public

import (
	"context"
	"fmt"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	util "beaver/utils/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetH5AuthCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// H5 免登获取 authCode（需在 WebView 环境中调用）
func NewGetH5AuthCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetH5AuthCodeLogic {
	return &GetH5AuthCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetH5AuthCodeLogic) GetH5AuthCode(req *types.GetH5AuthCodeReq) (resp *types.GetH5AuthCodeRes, err error) {
	// 1. 从 context 中获取当前用户 ID（由 AuthMiddleware 注入）
	userID, ok := l.ctx.Value("userID").(string)
	if !ok || userID == "" {
		logx.Error("H5免登失败：未获取到用户ID")
		return nil, fmt.Errorf("未登录")
	}

	// 2. 验证 appId 是否存在
	var app open_models.OpenApp
	if err := l.svcCtx.DB.Where("app_id = ?", req.AppID).First(&app).Error; err != nil {
		logx.Errorf("应用不存在: appId=%s, err=%v", req.AppID, err)
		return nil, fmt.Errorf("应用不存在")
	}

	// 3. 检查应用状态
	if app.Status != 1 {
		return nil, fmt.Errorf("应用未启用")
	}

	// 4. 生成 authCode（UUID）
	authCode := util.NewV4().String()

	// 5. 存储 authCode 到数据库（用于后续换取用户信息）
	// TODO: 需要创建 open_auth_code 表来存储 authCode
	// 这里暂时返回成功，实际应该存储到数据库
	logx.Infof("生成 H5 authCode: authCode=%s, appId=%s, userId=%s", authCode, req.AppID, userID)

	return &types.GetH5AuthCodeRes{
		AuthCode: authCode,
		ExpireIn: 300, // 5分钟 = 300秒
	}, nil
}
