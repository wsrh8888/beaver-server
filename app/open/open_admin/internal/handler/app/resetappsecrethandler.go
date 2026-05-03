package handler

import (
	"net/http"

	"beaver/app/open/open_admin/internal/logic/app"
	"beaver/app/open/open_admin/internal/svc"
	"beaver/app/open/open_admin/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ResetAppSecretHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ResetAppSecretReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := app.NewResetAppSecretLogic(r.Context(), svcCtx)
		resp, err := l.ResetAppSecret(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
