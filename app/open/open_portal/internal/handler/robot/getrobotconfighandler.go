package handler

import (
	logic "beaver/app/open/open_portal/internal/logic/robot"
	"beaver/app/open/open_portal/internal/svc"
	"beaver/app/open/open_portal/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetRobotConfigHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetRobotConfigReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewGetRobotConfigLogic(r.Context(), svcCtx)
		resp, err := l.GetRobotConfig(&req)
		response.Response(r, w, resp, err)
	}
}
