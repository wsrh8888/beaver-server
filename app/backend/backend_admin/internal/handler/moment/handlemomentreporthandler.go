package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/moment"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func HandleMomentReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HandleMomentReportReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.Status < 1 || req.Status > 2 {
			response.Response(r, w, nil, errors.New("无效的处理状态，状态值应为：1-已处理，2-已驳回"))
			return
		}

		l := logic.NewHandleMomentReportLogic(r.Context(), svcCtx)
		resp, err := l.HandleMomentReport(&req)
		response.Response(r, w, resp, err)
	}
}
