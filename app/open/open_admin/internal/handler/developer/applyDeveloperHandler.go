package handler

import (
	"net/http"

	"beaver/app/open/open_admin/internal/logic/developer"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 申请成为开发者
func ApplyDeveloperHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyDeveloperReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := developer.NewApplyDeveloperLogic(r.Context(), svcCtx)
		resp, err := l.ApplyDeveloper(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
