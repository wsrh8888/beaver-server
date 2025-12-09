package wsCommandConst

// 定义消息命令和类型
type Command string

// 消息主类型
const (
	// 聊天消息类
	CHAT_MESSAGE Command = "CHAT_MESSAGE"
	// 好友关系类
	FRIEND_OPERATION Command = "FRIEND_OPERATION"
	// 群组操作类
	GROUP_OPERATION Command = "GROUP_OPERATION"
	// 用户信息类
	USER_PROFILE Command = "USER_PROFILE"
	// 通知中心类
	NOTIFICATION Command = "NOTIFICATION"
	// 心跳消息类（应用级心跳）
	HEARTBEAT Command = "HEARTBEAT"
)
