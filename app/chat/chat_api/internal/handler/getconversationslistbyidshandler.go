package handler

import (
	"net/http"

	"beaver/app/chat/chat_api/internal/logic"
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量获取会话数据
func getConversationsListByIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetConversationsListByIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetConversationsListByIdsLogic(r.Context(), svcCtx)
		resp, err := l.GetConversationsListByIds(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
