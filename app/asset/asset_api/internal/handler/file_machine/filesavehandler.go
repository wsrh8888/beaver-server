package handler

import (
	logic "beaver/app/asset/asset_api/internal/logic/file_machine"
	"beaver/app/asset/asset_api/internal/svc"
	"beaver/app/asset/asset_api/internal/types"
	"beaver/common/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveFileReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewFileSaveLogic(r.Context(), svcCtx)
		resp, err := l.FileSave(&req)
		response.Response(r, w, resp, err)
	}
}
