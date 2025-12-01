package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/emoji"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func DeleteEmojiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteEmojiReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.UUID == "" {
			response.Response(r, w, nil, errors.New("表情UUID不能为空"))
			return
		}

		l := logic.NewDeleteEmojiLogic(r.Context(), svcCtx)
		resp, err := l.DeleteEmoji(&req)
		response.Response(r, w, resp, err)
	}
}
