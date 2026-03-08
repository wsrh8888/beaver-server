package handler

import (
	"beaver/app/chat/chat_api/internal/logic"
	"beaver/app/chat/chat_api/internal/svc"
	"beaver/app/chat/chat_api/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getForwardDetailsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetForwardDetailsReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		if req.RecordID == "" {
			response.Response(r, w, nil, errors.New("recordId不能为空"))
			return
		}

		l := logic.NewGetForwardDetailsLogic(r.Context(), svcCtx)
		resp, err := l.GetForwardDetails(&req)
		response.Response(r, w, resp, err)
	}
}
