package handler

import (
	"beaver/app/track/track_api/internal/logic"
	"beaver/app/track/track_api/internal/svc"
	"beaver/app/track/track_api/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func reportEventsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReportEventsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数验证
		if len(req.Events) == 0 {
			response.Response(r, w, nil, errors.New("events cannot be empty"))
			return
		}

		l := logic.NewReportEventsLogic(r.Context(), svcCtx)
		resp, err := l.ReportEvents(&req)
		response.Response(r, w, resp, err)
	}
}
