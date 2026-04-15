package webhookconst

// Webhook 事件类型常量
const (
	// Bot 相关事件
	EventBotMessageReceive = "bot.message.receive" // Bot 收到消息

	// 消息相关事件
	EventAfterSendMsg  = "after.send_msg"  // 消息发送后
	EventBeforeSendMsg = "before.send_msg" // 消息发送前(可拦截)
	EventMessageRecall = "message.recall"  // 消息撤回

	// 好友相关事件
	EventAfterAddFriend    = "after.add_friend"    // 加好友后
	EventAfterDeleteFriend = "after.delete_friend" // 删除好友后

	// 群组相关事件
	EventAfterJoinGroup  = "after.join_group"  // 加群后
	EventAfterLeaveGroup = "after.leave_group" // 退群后
	EventGroupDissolve   = "group.dissolve"    // 群组解散

	// 用户相关事件
	EventUserRegister = "user.register" // 用户注册
	EventUserLogin    = "user.login"    // 用户登录
)
