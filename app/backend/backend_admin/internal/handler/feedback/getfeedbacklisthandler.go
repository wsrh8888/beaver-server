package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/feedback"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetFeedbackListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFeedbackListReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetFeedbackListLogic(r.Context(), svcCtx)
		resp, err := l.GetFeedbackList(&req)
		response.Response(r, w, resp, err)
	}
}
