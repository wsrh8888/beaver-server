package open_models

// OpenAppScope 应用权限范围定义（对标飞书/钉钉）
type OpenAppScope string

const (
	// OAuth 相关权限
	ScopeOpenID      OpenAppScope = "openid"       // 获取用户 OpenID
	ScopeUserProfile OpenAppScope = "user_profile" // 获取用户基本信息（昵称、头像）
	ScopeUserPhone   OpenAppScope = "user_phone"   // 获取用户手机号
	ScopeUserEmail   OpenAppScope = "user_email"   // 获取用户邮箱

	// Bot 机器人权限
	ScopeBotMessageSend OpenAppScope = "bot.message.send"   // 以机器人身份发送消息
	ScopeBotMessageRecv OpenAppScope = "bot.message.recv"   // 接收消息事件
	ScopeBotStream      OpenAppScope = "bot.message.stream" // 流式消息（AI对话）

	// 消息权限
	ScopeMessageSend   OpenAppScope = "message.send"   // 发送消息（文本、图片、文件等）
	ScopeMessageRecall OpenAppScope = "message.recall" // 撤回消息
	ScopeMessageUpdate OpenAppScope = "message.update" // 更新消息（卡片）

	// 通讯录权限
	ScopeContactUserRead  OpenAppScope = "contact.user.read"  // 读取用户信息
	ScopeContactUserWrite OpenAppScope = "contact.user.write" // 管理用户（创建、更新、删除）
	ScopeContactDeptRead  OpenAppScope = "contact.dept.read"  // 读取部门信息
	ScopeContactDeptWrite OpenAppScope = "contact.dept.write" // 管理部门

	// 群组权限
	ScopeGroupRead   OpenAppScope = "group.read"    // 读取群组信息
	ScopeGroupWrite  OpenAppScope = "group.write"   // 管理群组（创建、解散、添加成员）
	ScopeGroupBotAdd OpenAppScope = "group.bot.add" // 邀请机器人入群

	// Webhook 权限
	ScopeWebhookEvent    OpenAppScope = "webhook.event"    // 接收事件推送
	ScopeWebhookIncoming OpenAppScope = "webhook.incoming" // Incoming Webhook（外部通知）
)

// ScopeDescription 权限描述映射
var ScopeDescription = map[OpenAppScope]string{
	ScopeOpenID:           "获取用户唯一标识",
	ScopeUserProfile:      "获取用户昵称、头像等基本信息",
	ScopeUserPhone:        "获取用户手机号",
	ScopeUserEmail:        "获取用户邮箱",
	ScopeBotMessageSend:   "以机器人身份发送消息",
	ScopeBotMessageRecv:   "接收用户发送给机器人的消息",
	ScopeBotStream:        "支持流式消息（适用于 AI 对话场景）",
	ScopeMessageSend:      "发送文本、图片、文件等消息",
	ScopeMessageRecall:    "撤回已发送的消息",
	ScopeMessageUpdate:    "更新消息内容（如动态更新卡片）",
	ScopeContactUserRead:  "读取用户详细信息",
	ScopeContactUserWrite: "创建、更新、删除用户",
	ScopeContactDeptRead:  "读取部门结构和信息",
	ScopeContactDeptWrite: "创建、更新、删除部门",
	ScopeGroupRead:        "读取群组信息和成员列表",
	ScopeGroupWrite:       "创建群组、解散群组、管理成员",
	ScopeGroupBotAdd:      "邀请机器人加入群组",
	ScopeWebhookEvent:     "接收消息、好友申请等事件推送",
	ScopeWebhookIncoming:  "接收外部系统（Jenkins/GitHub）的 Webhook 通知",
}

// DefaultScopes 创建应用时的默认权限（最小权限原则）
var DefaultScopes = []OpenAppScope{
	ScopeOpenID,
	ScopeUserProfile,
}

// AllScopes 所有可用权限
var AllScopes = []OpenAppScope{
	ScopeOpenID,
	ScopeUserProfile,
	ScopeUserPhone,
	ScopeUserEmail,
	ScopeBotMessageSend,
	ScopeBotMessageRecv,
	ScopeBotStream,
	ScopeMessageSend,
	ScopeMessageRecall,
	ScopeMessageUpdate,
	ScopeContactUserRead,
	ScopeContactUserWrite,
	ScopeContactDeptRead,
	ScopeContactDeptWrite,
	ScopeGroupRead,
	ScopeGroupWrite,
	ScopeGroupBotAdd,
	ScopeWebhookEvent,
	ScopeWebhookIncoming,
}
