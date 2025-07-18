package handler

import (
	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"beaver/common/validator"
	"beaver/utils/email"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func emailRegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmailRegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 验证邮箱格式
		if !email.IsValidEmail(req.Email) {
			response.Response(r, w, nil, errors.New("邮箱格式不正确"))
			return
		}

		// 验证验证码格式（6位纯数字）
		if !email.IsValidVerificationCode(req.Code) {
			response.Response(r, w, nil, errors.New("验证码格式不正确"))
			return
		}

		// 验证密码强度
		if !validator.IsValidPassword(req.Password) {
			response.Response(r, w, nil, errors.New("密码必须包含数字和字母,且长度至少8位"))
			return
		}

		l := logic.NewEmailRegisterLogic(r.Context(), svcCtx)
		resp, err := l.EmailRegister(&req)
		response.Response(r, w, resp, err, "注册成功")
	}
}
