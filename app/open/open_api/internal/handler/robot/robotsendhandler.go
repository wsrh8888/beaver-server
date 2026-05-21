package handler

import (
	logic "beaver/app/open/open_api/internal/logic/robot"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RobotSendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RobotSendReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewRobotSendLogic(r.Context(), svcCtx)
		resp, err := l.RobotSend(&req)
		response.Response(r, w, resp, err)
	}
}
