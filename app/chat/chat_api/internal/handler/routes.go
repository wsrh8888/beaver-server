// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"beaver/app/chat/chat_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/deleteRecentChat",
				Handler: deleteRecentHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/edit",
				Handler: editMessageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/forward",
				Handler: forwardMessageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/getChatHistory",
				Handler: chatHistoryHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/getConversationInfo",
				Handler: ConversationInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/chat/getRecentChatList",
				Handler: recentChatListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/pinnedChat",
				Handler: pinnedChatHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/recall",
				Handler: recallMessageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/chat/sendMsg",
				Handler: SendMsgHandler(serverCtx),
			},
		},
	)
}
