package handler

import (
	logic "beaver/app/backend/backend_admin/internal/logic/user"
	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/common/response"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetUserDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.UserID == "" {
			response.Response(r, w, nil, errors.New("用户ID不能为空"))
			return
		}

		l := logic.NewGetUserDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetUserDetail(&req)
		response.Response(r, w, resp, err)
	}
}
