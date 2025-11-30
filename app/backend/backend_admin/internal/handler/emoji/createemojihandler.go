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

func CreateEmojiHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateEmojiReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.FileKey == "" {
			response.Response(r, w, nil, errors.New("文件ID不能为空"))
			return
		}
		if req.Title == "" {
			response.Response(r, w, nil, errors.New("表情名称不能为空"))
			return
		}
		if req.AuthorID == "" {
			response.Response(r, w, nil, errors.New("创建者ID不能为空"))
			return
		}

		l := logic.NewCreateEmojiLogic(r.Context(), svcCtx)
		resp, err := l.CreateEmoji(&req)
		response.Response(r, w, resp, err)
	}
}
