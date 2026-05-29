package track

import (
	"errors"
	"net/http"

	tracklogic "beaver/app/platform/platform_api/internal/logic/track"
	"beaver/app/platform/platform_api/internal/svc"
	"beaver/app/platform/platform_api/internal/types"
	"beaver/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ReportEventsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReportEventsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if len(req.Events) == 0 {
			response.Response(r, w, nil, errors.New("events cannot be empty"))
			return
		}

		l := tracklogic.NewReportEventsLogic(r.Context(), svcCtx)
		resp, err := l.ReportEvents(&req)
		response.Response(r, w, resp, err)
	}
}
