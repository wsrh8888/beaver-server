package handler

import (
	moderationLogic "beaver/app/backend/backend_admin/internal/logic/moderation"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RejectContentReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RejectContentReportReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := moderationLogic.NewRejectContentReportLogic(r.Context(), svcCtx)
		resp, err := l.RejectContentReport(&req)
		response.Response(r, w, resp, err)
	}
}
