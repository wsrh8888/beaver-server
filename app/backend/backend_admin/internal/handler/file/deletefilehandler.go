package handler

import (
	"beaver/app/backend/backend_admin/internal/logic/file"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteFileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewDeleteFileLogic(r.Context(), svcCtx)
		resp, err := l.DeleteFile(&req)
		response.Response(r, w, resp, err)
	}
}
