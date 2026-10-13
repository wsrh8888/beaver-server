/*
 * 共享数据类型：测试账号与好友关系，供 setup 脚本和压测脚本共用。
 *
 * 纯类型定义，独立成 package types，不依赖任何业务逻辑。
 */

package types

// Account 测试账号，序列化到 accounts.json 供所有压测脚本共用
type Account struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	UserID   string `json:"userId,omitempty"`
	Token    string `json:"token,omitempty"`
}

// TestPassword 测试账号统一密码
const TestPassword = "Test@123456"

// Friendship 一对好友关系 + 自动产生的会话ID，供 chat 压测取 conversationId。
// 序列化到 friends.json，由 setup/friends 产出。
type Friendship struct {
	UserID         string `json:"userId"`
	FriendID       string `json:"friendId"`
	ConversationID string `json:"conversationId"`
}
