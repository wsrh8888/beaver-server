package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/emoji"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func AddEmojiToPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AddEmojiToPackageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewAddEmojiToPackageLogic(r.Context(), svcCtx)
		resp, err := l.AddEmojiToPackage(&req)
		response.Response(r, w, resp, err)
	}
}
