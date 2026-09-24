/*
 * 公共账号类型，供 seed 脚本和压测脚本共用
 */

package lib

// Account 测试账号，序列化到 accounts.json 供所有压测脚本共用
type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   string `json:"userId,omitempty"`
	Token    string `json:"token,omitempty"`
}

// TestPassword 测试账号统一密码
const TestPassword = "Test@123456"
