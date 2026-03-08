package handler

import (
	"beaver/app/call/call_api/internal/logic"
	"beaver/app/call/call_api/internal/svc"
	"beaver/app/call/call_api/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func StartCallHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StartCallReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.ConversationId == "" {
			response.Response(r, w, nil, errors.New("会话ID不能为空"))
			return
		}
		if req.CallType != 1 && req.CallType != 2 {
			response.Response(r, w, nil, errors.New("通话类型不合法"))
			return
		}
		if req.CallMode != 1 && req.CallMode != 2 {
			response.Response(r, w, nil, errors.New("通话模式不合法"))
			return
		}

		l := logic.NewStartCallLogic(r.Context(), svcCtx)
		resp, err := l.StartCall(&req)
		response.Response(r, w, resp, err)
	}
}
