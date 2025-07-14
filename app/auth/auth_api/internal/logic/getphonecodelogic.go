package logic

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

func (l *GetPhoneCodeLogic) GetPhoneCode(req *types.GetPhoneCodeReq) (resp *types.GetPhoneCodeRes, err error) {
	// 生成6位数字验证码
	code := email.GenerateCode()

	// 发送手机验证码（这里需要集成短信服务）
	err = l.sendSMS(req.Phone, code, req.Type)
	if err != nil {
		logx.Errorf("发送短信失败: %v", err)
		return nil, errors.New("发送验证码失败，请稍后重试")
	}

	// 存储验证码到Redis（5分钟有效期）
	codeKey := fmt.Sprintf("phone_code_%s_%s", req.Phone, req.Type)
	err = l.svcCtx.Redis.Set(codeKey, code, 5*time.Minute).Err()
	if err != nil {
		logx.Errorf("存储验证码失败: %v", err)
		return nil, errors.New("服务内部异常")
	}

	// 设置发送频率限制（60秒）
	rateLimitKey := fmt.Sprintf("phone_rate_limit_%s", req.Phone)
	err = l.svcCtx.Redis.Set(rateLimitKey, "1", 60*time.Second).Err()
	if err != nil {
		logx.Errorf("设置发送频率限制失败: %v", err)
	}

	return &types.GetPhoneCodeRes{
		Message: "验证码发送成功",
	}, nil
}

// 发送短信验证码
func (l *GetPhoneCodeLogic) sendSMS(phone, code, codeType string) error {
	// TODO: 集成短信服务商（如阿里云、腾讯云等）
	// 这里暂时返回成功，实际项目中需要集成真实的短信服务

	logx.Infof("发送短信验证码到 %s: %s (类型: %s)", phone, code, codeType)

	// 示例：调用短信服务API
	// smsClient := sms.NewClient(accessKeyId, accessKeySecret)
	// request := &sms.SendSmsRequest{
	//     PhoneNumbers:  phone,
	//     SignName:      "海狸IM",
	//     TemplateCode:  getSMSTemplate(codeType),
	//     TemplateParam: fmt.Sprintf(`{"code":"%s"}`, code),
	// }
	// _, err := smsClient.SendSms(request)
	// return err

	return nil
}

// 获取短信模板代码
func getSMSTemplate(codeType string) string {
	switch codeType {
	case "register":
		return "SMS_123456789" // 注册模板
	case "reset":
		return "SMS_987654321" // 重置密码模板
	case "login":
		return "SMS_456789123" // 登录模板
	default:
		return "SMS_123456789" // 默认模板
	}
}
