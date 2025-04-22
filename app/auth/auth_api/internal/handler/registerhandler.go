package handler

import (
	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"beaver/common/validator"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func registerHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 验证手机号格式
		if !validator.IsValidPhone(req.Phone) {
			response.Response(r, w, nil, errors.New("手机号格式不正确"))
			return
		}

		// 验证密码强度
		if !validator.IsValidPassword(req.Password) {
			response.Response(r, w, nil, errors.New("密码必须包含数字和字母,且长度至少8位"))
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		response.Response(r, w, resp, err, "注册成功")
	}
}
