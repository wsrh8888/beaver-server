package handler

import (
	moderationLogic "beaver/app/backend/backend_admin/internal/logic/moderation"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func EscalateContentReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EscalateContentReportReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := moderationLogic.NewEscalateContentReportLogic(r.Context(), svcCtx)
		resp, err := l.EscalateContentReport(&req)
		response.Response(r, w, resp, err)
	}
}
