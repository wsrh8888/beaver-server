package wsTypeConst

type Type string

// send（发起消息给服务端）
// receive（发给对方设备）
// sync（发给自己其他设备进行记录同步）

const (
	PrivateMessageSend Type = "private_message_send" // 客户端->服务端 私聊消息发送
	GroupMessageSend   Type = "group_message_send"   // 客户端->服务端 群聊消息发送
	// --------------------------------------------------------
	// --------------------------------------------------------

	// 会话信息同步
	ChatConversationMetaReceive    Type = "chat_conversation_meta_receive"    //  服务端->客户端 会话信息同步
	ChatUserConversationReceive    Type = "chat_user_conversation_receive"    //  服务端->客户端 用户会话信息同步
	ChatConversationMessageReceive Type = "chat_conversation_message_receive" //  服务端->客户端 会话消息同步
)
const (
	// -------------------------------------------------------------------------------------
	FriendReceive       Type = "friend_receive"        // 服务端->客户端 好友信息同步
	FriendVerifyReceive Type = "friend_verify_receive" // 服务端->客户端 好友验证信息同步
)

// -------------------------------------------------------------------------------------

const (
	GroupReceive            Type = "group_receive"              // 服务端->客户端 群组信息同步
	GroupJoinRequestReceive Type = "group_join_request_receive" // 服务端->客户端 群成员添加请求
	GroupMemberReceive      Type = "group_member_receive"       // 服务端->客户端 群成员变动（加入，离开、被踢出等）

)

const (
	// --------------------------------------------------------
	UserReceive Type = "user_receive" // 服务端->客户端 用户信息同步
)

const (
	// 通知中心
	NotificationReceive Type = "notification_receive" // 服务端->客户端 通知提醒
)
