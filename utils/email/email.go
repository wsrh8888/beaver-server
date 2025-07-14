package email

import (
	"fmt"
	"regexp"
	"strings"

	utils "beaver/utils/rand"
)

// 验证邮箱格式
func IsValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// 验证验证码类型
func IsValidCodeType(codeType string) bool {
	validTypes := []string{"register", "login", "reset_password", "update_email"}
	for _, t := range validTypes {
		if t == codeType {
			return true
		}
	}
	return false
}

// 验证验证码格式（6位纯数字）
func IsValidVerificationCode(code string) bool {
	if len(code) != 6 {
		return false
	}
	matched, _ := regexp.MatchString(`^\d{6}$`, code)
	return matched
}

// 生成6位数字验证码
func GenerateCode() string {
	return utils.GenerateNumericCode(6)
}

// 根据邮箱域名获取邮箱服务商
func GetEmailProvider(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}

	domain := strings.ToLower(parts[1])

	switch domain {
	case "qq.com":
		return "QQ"
	default:
		return ""
	}
}

// 获取邮件主题
func GetEmailSubject(codeType string) string {
	switch codeType {
	case "register":
		return "海狸IM - 注册验证码"
	case "login":
		return "海狸IM - 登录验证码"
	case "reset_password":
		return "海狸IM - 找回密码验证码"
	case "update_email":
		return "海狸IM - 修改邮箱验证码"
	default:
		return "海狸IM - 验证码"
	}
}

// 获取邮件内容
func GetEmailBody(code, codeType string) string {
	action := ""
	switch codeType {
	case "register":
		action = "注册"
	case "login":
		action = "登录"
	case "reset_password":
		action = "找回密码"
	case "update_email":
		action = "修改邮箱"
	default:
		action = "验证"
	}

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>海狸IM验证码</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
    <div style="max-width: 600px; margin: 0 auto; padding: 20px;">
        <h2 style="color: #2c3e50;">海狸IM - %s验证码</h2>
        <p>您好！</p>
        <p>您正在进行%s操作，验证码如下：</p>
        <div style="background-color: #f8f9fa; padding: 15px; border-radius: 5px; text-align: center; margin: 20px 0;">
            <h1 style="color: #e74c3c; font-size: 32px; margin: 0; letter-spacing: 5px;">%s</h1>
        </div>
        <p><strong>验证码有效期：5分钟</strong></p>
        <p>如果这不是您的操作，请忽略此邮件。</p>
        <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
        <p style="color: #7f8c8d; font-size: 12px;">
            此邮件由海狸IM系统自动发送，请勿回复。<br>
            如有疑问，请联系客服。
        </p>
    </div>
</body>
</html>
`, action, action, code)
}
