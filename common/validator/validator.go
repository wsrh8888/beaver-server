package validator

import (
	"regexp"
	"errors"
)

// 验证手机号格式
func IsValidPhone(phone string) bool {
	pattern := `^1[3-9]\d{9}$`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(phone)
}

// 验证密码强度
func IsValidPassword(password string) bool {
	// 密码长度至少8位
	if len(password) < 8 {
		return false
	}
	// 必须包含数字和字母
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	return hasNumber && hasLetter
}

// 验证登录参数
func ValidateLoginParams(phone, password string) error {
	if !IsValidPhone(phone) {
		return errors.New("手机号格式不正确")
	}
	if password == "" {
		return errors.New("密码不能为空")
	}
	return nil
} 