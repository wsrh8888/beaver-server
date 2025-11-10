package handler

import (
	"net/http"

	"beaver/app/moment/moment_api/internal/logic"
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量获取评论数据（用于数据同步）
func GetCommentsBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetCommentsBatchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetCommentsBatchLogic(r.Context(), svcCtx)
		resp, err := l.GetCommentsBatch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
