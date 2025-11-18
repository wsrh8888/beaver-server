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

func CreateEmojiPackageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateEmojiPackageReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.Title == "" {
			response.Response(r, w, nil, errors.New("表情包名称不能为空"))
			return
		}
		if req.UserID == "" {
			response.Response(r, w, nil, errors.New("创建者ID不能为空"))
			return
		}
		if req.Type == "" {
			response.Response(r, w, nil, errors.New("表情包类型不能为空"))
			return
		}

		l := logic.NewCreateEmojiPackageLogic(r.Context(), svcCtx)
		resp, err := l.CreateEmojiPackage(&req)
		response.Response(r, w, resp, err)
	}
}
