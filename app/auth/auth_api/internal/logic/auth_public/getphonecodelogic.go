package auth_public

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/utils/email"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/zeromicro/go-zero/core/logx"
)


type GetPhoneCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewGetPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPhoneCodeLogic {
	return &GetPhoneCodeLogic{
		ctx:    ctx,
		logger: logger.New("get_phone_code"),
		svcCtx: svcCtx,
	}
}

func (l *GetPhoneCodeLogic) GetPhoneCode(req *types.GetPhoneCodeReq) (*types.GetPhoneCodeRes, error) {
	rateLimitKey := fmt.Sprintf("phone_rate_limit_%s", req.Phone)
	exists, err := l.svcCtx.Redis.Exists(rateLimitKey).Result()
	if err != nil {
		logx.Errorf("检查短信发送频率限制失败: %v", err)
		return nil, errors.New("服务内部异常")
	}
	if exists > 0 {
		return nil, errors.New("发送过于频繁，请60秒后再试")
	}

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

	_ = l.svcCtx.Redis.Set(rateLimitKey, "1", 60*time.Second).Err()

	l.logger.Info(model.LogMsg{
		Text: "短信验证码已发送",
		Data: map[string]interface{}{"codeType": req.Type},
	})

	return &types.GetPhoneCodeRes{}, nil
}

func (l *GetPhoneCodeLogic) sendSMS(phone, code, codeType string) error {
	return nil
}