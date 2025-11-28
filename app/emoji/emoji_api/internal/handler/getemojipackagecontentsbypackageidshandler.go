package handler

import (
	"net/http"

	"beaver/app/emoji/emoji_api/internal/logic"
	"beaver/app/emoji/emoji_api/internal/svc"
	"beaver/app/emoji/emoji_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量获取表情包内容详情（同步用）
func GetEmojiPackageContentsByPackageIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetEmojiPackageContentsByPackageIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetEmojiPackageContentsByPackageIdsLogic(r.Context(), svcCtx)
		resp, err := l.GetEmojiPackageContentsByPackageIds(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
