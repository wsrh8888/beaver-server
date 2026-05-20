package websocket

// Command WebSocket 命令类型
type Command string

// 业务命令（携带完整 content/data 层）
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
	// 表情中心类
	EMOJI Command = "EMOJI"
	// 音视频通话类
	CALL Command = "CALL"
)

// 控制帧（无 content/data 层，仅携带 messageId / timestamp）
const (
	// PING 心跳探测，发送方期待对端回复 PONG
	PING Command = "PING"
	// PONG 心跳回应，收到 PING 后立即回复
	PONG Command = "PONG"
	// ACK 消息回执，服务端收到客户端消息后立即发送，表示"已收到"，无失败状态
	ACK Command = "ACK"
)
