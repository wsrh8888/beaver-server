package auth_public

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/email"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPhoneCodeLogic {
	return &GetPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPhoneCodeLogic) GetPhoneCode(req *types.GetPhoneCodeReq) (*types.GetPhoneCodeRes, error) {
	code := email.GenerateCode()
	if err := l.sendSMS(req.Phone, code, req.Type); err != nil {
		logx.Errorf("发送短信失败: %v", err)
		return nil, errors.New("发送验证码失败，请稍后重试")
	}

	codeKey := fmt.Sprintf("phone_code_%s_%s", req.Phone, req.Type)
	if err := l.svcCtx.Redis.Set(codeKey, code, 5*time.Minute).Err(); err != nil {
		logx.Errorf("存储验证码失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	rateLimitKey := fmt.Sprintf("phone_rate_limit_%s", req.Phone)
	_ = l.svcCtx.Redis.Set(rateLimitKey, "1", 60*time.Second).Err()

	return &types.GetPhoneCodeRes{}, nil
}

func (l *GetPhoneCodeLogic) sendSMS(phone, code, codeType string) error {
	logx.Infof("发送短信验证码到 %s: %s (类型: %s)", phone, code, codeType)
	return nil
}