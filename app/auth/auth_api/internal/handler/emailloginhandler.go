package handler

import (
	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"beaver/utils/device"
	"beaver/utils/email"
	"errors"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func emailLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EmailLoginReq
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

		// 验证设备ID
		if !device.IsValidDeviceID(req.DeviceID) {
			response.Response(r, w, nil, errors.New("设备ID格式不正确"))
			return
		}

		l := logic.NewEmailLoginLogic(r.Context(), svcCtx)
		resp, err := l.EmailLogin(&req)
		response.Response(r, w, resp, err)
	}
}
