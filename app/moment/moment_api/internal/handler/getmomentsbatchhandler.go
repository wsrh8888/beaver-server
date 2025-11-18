package handler

import (
	"net/http"

	"beaver/app/moment/moment_api/internal/logic"
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 批量获取动态数据（用于数据同步）
func GetMomentsBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetMomentsBatchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetMomentsBatchLogic(r.Context(), svcCtx)
		resp, err := l.GetMomentsBatch(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
