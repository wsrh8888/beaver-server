package handler

import (
	"encoding/json"
	"net/http"

	"beaver/app/datasync/datasync_api/internal/svc"
	chat_rpc "beaver/app/message/message_rpc/types/chat_rpc"
	"beaver/app/relation/relation_rpc/types"
	accttypes "beaver/app/account/account_rpc/types"
	assettypes "beaver/app/asset/asset_rpc/types"
	platform_rpc "beaver/app/platform/platform_rpc/types/platform_rpc"
	notification_rpc "beaver/app/notification/notification_rpc/types/notification_rpc"
	"beaver/common/response"
	"github.com/zeromicro/go-zero/rest"
)

// RegisterS2SHandlers 把内部编排面的业务端点挂上 S2S Bearer 验签中间件。
//
// 收费模块（beaver-social 等）经 client_credentials 换票后带 Bearer 调用；端点零转换透传
// 本面 gRPC 客户端（入参/出参均为 proto 类型，与调用方逐字段对齐），机器对机器，无终端用户身份。
func RegisterS2SHandlers(server *rest.Server, ctx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{ctx.InternalToken},
			[]rest.Route{
				// 用户信息批量查询（圈子/动态展示发帖人、评论人、点赞人等）
				{Method: http.MethodPost, Path: "/internal/User.ListInfo", Handler: UserListInfoHandler(ctx)},
				// 单个用户信息查询（通话发起人/被邀请人资料展示）
				{Method: http.MethodPost, Path: "/internal/User.UserInfo", Handler: UserInfoHandler(ctx)},
				// 密码校验（开放平台账号密码登录）
				{Method: http.MethodPost, Path: "/internal/Auth.VerifyPassword", Handler: AuthVerifyPasswordHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Auth.CreateCredential", Handler: AuthCreateCredentialHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Auth.ListUserDevices", Handler: AuthListUserDevicesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Auth.ResetPassword", Handler: AuthResetPasswordHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.BatchDeleteFiles", Handler: FileBatchDeleteFilesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.DeleteFile", Handler: FileDeleteFileHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.GetFileById", Handler: FileGetFileByIdHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.GetFileByMd5", Handler: FileGetFileByMd5Handler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.GetFileDetail", Handler: FileGetFileDetailHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.GetQiniuUploadToken", Handler: FileGetQiniuUploadTokenHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.ListFiles", Handler: FileListFilesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/File.SaveFile", Handler: FileSaveFileHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Emoji.ListEmojiCollects", Handler: EmojiListEmojiCollectsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Emoji.ListEmojiPackages", Handler: EmojiListEmojiPackagesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Emoji.ListEmojis", Handler: EmojiListEmojisHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Emoji.SaveEmoji", Handler: EmojiSaveEmojiHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Emoji.SaveEmojiPackage", Handler: EmojiSaveEmojiPackageHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Emoji.UpdateEmojiPackageContent", Handler: EmojiUpdateEmojiPackageContentHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.AdminQueryLogs", Handler: PlatformAdminQueryLogsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.CreateApp", Handler: PlatformCreateAppHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.CreateArchitecture", Handler: PlatformCreateArchitectureHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.CreateVersion", Handler: PlatformCreateVersionHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.CreateWorkbenchApp", Handler: PlatformCreateWorkbenchAppHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.DeleteFeedback", Handler: PlatformDeleteFeedbackHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.DeleteVersion", Handler: PlatformDeleteVersionHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.DeleteWorkbenchApp", Handler: PlatformDeleteWorkbenchAppHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.GetContentReport", Handler: PlatformGetContentReportHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.GetFeedback", Handler: PlatformGetFeedbackHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.GetWorkbenchApp", Handler: PlatformGetWorkbenchAppHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.HandleFeedback", Handler: PlatformHandleFeedbackHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListAppVersions", Handler: PlatformListAppVersionsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListApps", Handler: PlatformListAppsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListArchitectures", Handler: PlatformListArchitecturesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListContentReports", Handler: PlatformListContentReportsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListFeedback", Handler: PlatformListFeedbackHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListReleasePolicies", Handler: PlatformListReleasePoliciesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListVersions", Handler: PlatformListVersionsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.ListWorkbenchApps", Handler: PlatformListWorkbenchAppsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.UpdateArchitecture", Handler: PlatformUpdateArchitectureHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.UpdateContentReports", Handler: PlatformUpdateContentReportsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.UpdateWorkbenchApp", Handler: PlatformUpdateWorkbenchAppHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Platform.UpsertReleasePolicy", Handler: PlatformUpsertReleasePolicyHandler(ctx)},
				// 会话初始化（建圈/加圈/邀请后同步初始化 IM 会话）
				{Method: http.MethodPost, Path: "/internal/Chat.InitializeConversation", Handler: ChatInitializeConversationHandler(ctx)},
				// 会话内通知消息（xxx 创建了圈子、xxx 加入了圈子等系统通知）
				{Method: http.MethodPost, Path: "/internal/Chat.SendNotificationMessage", Handler: ChatSendNotificationMessageHandler(ctx)},
				// 消息发送（通话挂断后往会话投递通话消息）
				{Method: http.MethodPost, Path: "/internal/Chat.SendMsg", Handler: ChatSendMsgHandler(ctx)},
				// 通知事件推送（动态评论/点赞提醒等）
				{Method: http.MethodPost, Path: "/internal/Notification.PushEvent", Handler: NotificationPushEventHandler(ctx)},
				// 好友 ID 列表（动态可见性按好友圈分发）
				{Method: http.MethodPost, Path: "/internal/Friend.GetFriendIds", Handler: FriendGetFriendIdsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.UserCreate", Handler: UserUserCreateHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.SearchUser", Handler: UserSearchUserHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.UserUpdateDisplay", Handler: UserUserUpdateDisplayHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.ListUsers", Handler: UserListUsersHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.UpdateUsers", Handler: UserUpdateUsersHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.DeleteUsers", Handler: UserDeleteUsersHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/User.UpdateUsersStatus", Handler: UserUpdateUsersStatusHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Chat.GetChatMessage", Handler: ChatGetChatMessageHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Chat.ListChatMessages", Handler: ChatListChatMessagesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Chat.UpdateChatMessages", Handler: ChatUpdateChatMessagesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Chat.ListConversations", Handler: ChatListConversationsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Friend.ListFriends", Handler: FriendListFriendsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Friend.UpdateFriends", Handler: FriendUpdateFriendsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Friend.ListFriendVerifies", Handler: FriendListFriendVerifiesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Friend.UpdateFriendVerifies", Handler: FriendUpdateFriendVerifiesHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Friend.ListFriendBlocks", Handler: FriendListFriendBlocksHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Friend.UpdateFriendBlocks", Handler: FriendUpdateFriendBlocksHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.GetUserGroupIDs", Handler: GroupGetUserGroupIDsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.GetGroupsListByIds", Handler: GroupGetGroupsListByIdsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.AddGroupMember", Handler: GroupAddGroupMemberHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.RemoveGroupMember", Handler: GroupRemoveGroupMemberHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.ListGroups", Handler: GroupListGroupsHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.UpdateGroup", Handler: GroupUpdateGroupHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.ListGroupMembers", Handler: GroupListGroupMembersHandler(ctx)},
				{Method: http.MethodPost, Path: "/internal/Group.UpdateGroupMember", Handler: GroupUpdateGroupMemberHandler(ctx)},
			}...,
		),
	)
}

func AuthVerifyPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.VerifyPasswordReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.AuthRpc.VerifyPassword(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.UserInfoReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.UserInfo(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserListInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.UserListInfoReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.UserListInfo(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatInitializeConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.InitializeConversationReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.InitializeConversation(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatSendMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.SendMsgReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.SendMsg(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatSendNotificationMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.SendNotificationMessageReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.SendNotificationMessage(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func NotificationPushEventHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req notification_rpc.PushEventReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.NotificationRpc.PushEvent(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendGetFriendIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFriendIdsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.GetFriendIds(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserUserCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.UserCreateReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.UserCreate(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserSearchUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.SearchUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.SearchUser(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserUserUpdateDisplayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.UserUpdateDisplayReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.UserUpdateDisplay(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserListUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.ListUsersReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.ListUsers(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserUpdateUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.UpdateUsersReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.UpdateUsers(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserDeleteUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.DeleteUsersReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.DeleteUsers(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func UserUpdateUsersStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.UpdateUsersStatusReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.UserRpc.UpdateUsersStatus(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatGetChatMessageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.GetChatMessageReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.GetChatMessage(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatListChatMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.ListChatMessagesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.ListChatMessages(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatUpdateChatMessagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.UpdateChatMessagesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.UpdateChatMessages(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func ChatListConversationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req chat_rpc.ListConversationsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.ChatRpc.ListConversations(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendListFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListFriendsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.ListFriends(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendUpdateFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFriendsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.UpdateFriends(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendListFriendVerifiesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListFriendVerifiesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.ListFriendVerifies(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendUpdateFriendVerifiesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFriendVerifiesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.UpdateFriendVerifies(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendListFriendBlocksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListFriendBlocksReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.ListFriendBlocks(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FriendUpdateFriendBlocksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateFriendBlocksReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FriendRpc.UpdateFriendBlocks(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupGetUserGroupIDsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserGroupIDsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.GetUserGroupIDs(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupGetGroupsListByIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGroupsListByIdsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.GetGroupsListByIds(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupAddGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddGroupMemberReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.AddGroupMember(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupRemoveGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveGroupMemberReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.RemoveGroupMember(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupListGroupsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListGroupsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.ListGroups(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupUpdateGroupHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGroupReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.UpdateGroup(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupListGroupMembersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListGroupMembersReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.ListGroupMembers(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func GroupUpdateGroupMemberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateGroupMemberReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.GroupRpc.UpdateGroupMember(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func AuthCreateCredentialHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.CreateCredentialReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.AuthRpc.CreateCredential(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func AuthListUserDevicesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.ListUserDevicesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.AuthRpc.ListUserDevices(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func AuthResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req accttypes.ResetPasswordReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.AuthRpc.ResetPassword(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileBatchDeleteFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.BatchDeleteFilesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.BatchDeleteFiles(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileDeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.DeleteFileReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.DeleteFile(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileGetFileByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.GetFileByIdReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.GetFileById(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileGetFileByMd5Handler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.GetFileByMd5Req
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.GetFileByMd5(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileGetFileDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.GetFileDetailReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.GetFileDetail(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileGetQiniuUploadTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.GetQiniuUploadTokenReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.GetQiniuUploadToken(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileListFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.ListFilesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.ListFiles(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func FileSaveFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.SaveFileReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.FileRpc.SaveFile(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func EmojiListEmojiCollectsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.ListEmojiCollectsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.EmojiRpc.ListEmojiCollects(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func EmojiListEmojiPackagesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.ListEmojiPackagesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.EmojiRpc.ListEmojiPackages(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func EmojiListEmojisHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.ListEmojisReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.EmojiRpc.ListEmojis(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func EmojiSaveEmojiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.SaveEmojiReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.EmojiRpc.SaveEmoji(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func EmojiSaveEmojiPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.SaveEmojiPackageReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.EmojiRpc.SaveEmojiPackage(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func EmojiUpdateEmojiPackageContentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req assettypes.UpdateEmojiPackageContentReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.EmojiRpc.UpdateEmojiPackageContent(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformAdminQueryLogsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.AdminQueryLogsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.AdminQueryLogs(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}






func PlatformCreateAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.CreateAppReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.CreateApp(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformCreateArchitectureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.CreateArchitectureReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.CreateArchitecture(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}



func PlatformCreateVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.CreateVersionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.CreateVersion(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformCreateWorkbenchAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.CreateWorkbenchAppReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.CreateWorkbenchApp(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}








func PlatformDeleteFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.DeleteFeedbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.DeleteFeedback(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformDeleteVersionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.DeleteVersionReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.DeleteVersion(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformDeleteWorkbenchAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.DeleteWorkbenchAppReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.DeleteWorkbenchApp(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}







func PlatformGetContentReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.GetContentReportReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.GetContentReport(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}


func PlatformGetFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.GetFeedbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.GetFeedback(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformGetWorkbenchAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.GetWorkbenchAppReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.GetWorkbenchApp(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformHandleFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.HandleFeedbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.HandleFeedback(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}






func PlatformListAppVersionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListAppVersionsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListAppVersions(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformListAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListAppsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListApps(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformListArchitecturesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListArchitecturesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListArchitectures(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}


func PlatformListContentReportsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListContentReportsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListContentReports(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}


func PlatformListFeedbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListFeedbackReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListFeedback(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformListReleasePoliciesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListReleasePoliciesReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListReleasePolicies(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformListVersionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListVersionsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListVersions(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformListWorkbenchAppsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.ListWorkbenchAppsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.ListWorkbenchApps(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}









func PlatformUpdateArchitectureHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.UpdateArchitectureReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.UpdateArchitecture(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}


func PlatformUpdateContentReportsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.UpdateContentReportsReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.UpdateContentReports(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}


func PlatformUpdateWorkbenchAppHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.UpdateWorkbenchAppReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.UpdateWorkbenchApp(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}

func PlatformUpsertReleasePolicyHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req platform_rpc.UpsertReleasePolicyReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		resp, err := svcCtx.PlatformRpc.UpsertReleasePolicy(r.Context(), &req)
		response.Response(r, w, resp, err)
	}
}
