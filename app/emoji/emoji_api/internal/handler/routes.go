// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"beaver/app/emoji/emoji_api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/add",
				Handler: AddEmojiHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/favoriteEmoji",
				Handler: UpdateFavoriteEmojiHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/favoriteList",
				Handler: GetEmojisListHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/favoritePackageList",
				Handler: GetUserFavoritePackagesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageAddEmoji",
				Handler: AddEmojiToPackageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageBatchAdd",
				Handler: BatchAddEmojiToPackageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageCreate",
				Handler: CreateEmojiPackageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageDeleteEmoji",
				Handler: DeleteEmojiFromPackageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageFavorite",
				Handler: UpdateFavoriteEmojiPackageHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageInfo",
				Handler: GetEmojiPackageDetailHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/emoji/packageList",
				Handler: GetEmojiPackagesHandler(serverCtx),
			},
		},
	)
}
