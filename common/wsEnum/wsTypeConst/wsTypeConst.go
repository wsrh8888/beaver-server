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
	FriendAddRequest     Type = "friend_add_request"     // 客户端->服务端 添加好友请求
	FriendAccept         Type = "friend_accept"          // 客户端->服务端 接受好友请求
	FriendReject         Type = "friend_reject"          // 客户端->服务端 拒绝好友请求
	FriendDelete         Type = "friend_delete"          // 服务端->客户端 删除好友
	FriendRequestReceive Type = "friend_request_receive" //  服务端->客户端 收到好友请求
	FriendAddSuccess     Type = "friend_add_success"     //  服务端->客户端 好友添加成功
)

const (
	GroupCreate      Type = "group_create"       // 客户端->服务端 创建群组
	GroupInvite      Type = "group_invite"       // 客户端->服务端 邀请入群
	GroupJoinRequest Type = "group_join_request" // 客户端->服务端 申请入群
	GroupQuit        Type = "group_quit"         // 客户端->服务端 退出群组

	GroupInviteReceive Type = "group_invite_receive" // 服务端->客户端 群聊消息接收
	GroupJoinApprove   Type = "group_join_approve"   // 服务端->客户端 群成员添加请求
	GroupMemberUpdate  Type = "group_member_update"  // 服务端->客户端 群成员变动（加入，离开、被踢出等）
	MessageGroupCreate Type = "message_group_create" // 服务端->客户端 创建群聊
	GroupUpdate        Type = "group_update"         // 服务端->客户端 群组更新（包含群信息更新、群主转让等）
)

const (
	ProfileUpdate       Type = "profile_update"        // 客户端->服务端 更新个人信息
	ProfileChangeNotify Type = "profile_change_notify" // 服务端->客户端 他人资料变更通知
)
