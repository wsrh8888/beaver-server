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

func CreateUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateUserReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 参数校验
		if req.Nickname == "" || req.Password == "" || req.Email == "" {
			response.Response(r, w, nil, errors.New("昵称、密码和邮箱不能为空"))
			return
		}

		l := logic.NewCreateUserLogic(r.Context(), svcCtx)
		resp, err := l.CreateUser(&req)
		response.Response(r, w, resp, err)
	}
}
