package oauth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDesktopQuickLoginSignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 PC 端快捷登录签名
func NewGetDesktopQuickLoginSignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDesktopQuickLoginSignLogic {
	return &GetDesktopQuickLoginSignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDesktopQuickLoginSignLogic) GetDesktopQuickLoginSign(req *types.GetDesktopQuickLoginSignReq) (resp *types.GetDesktopQuickLoginSignRes, err error) {
	// 1. token 已由中间件验证，从 context 取 appID
	appID, _ := l.ctx.Value("appID").(string)
	_ = appID

	// 2. 验证 token 有效性（调用 OAuth RPC）
	validateResp, err := l.svcCtx.OAuthRpc.ValidateToken(l.ctx, nil) // TODO: 需要传入 token
	if err != nil || !validateResp.Valid {
		return nil, errors.New("无效的访问令牌")
	}

	// 3. 生成时间戳（毫秒）
	timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// 4. 使用 appSecret 计算签名（HmacSHA256）
	// TODO: 需要从应用配置中获取 appSecret
	appSecret := "your_app_secret" // 临时硬编码，实际应从数据库获取

	mac := hmac.New(sha256.New, []byte(appSecret))
	mac.Write([]byte(timestamp))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return &types.GetDesktopQuickLoginSignRes{
		QuickSign: signature,
	}, nil
}
