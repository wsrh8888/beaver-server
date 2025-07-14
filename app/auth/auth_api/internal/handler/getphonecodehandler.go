package handler

import (
	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"beaver/common/validator"
	"beaver/utils/email"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func getPhoneCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetPhoneCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 验证手机号格式
		if !validator.IsValidPhone(req.Phone) {
			response.Response(r, w, nil, errors.New("手机号格式不正确"))
			return
		}

		// 验证验证码类型
		if !email.IsValidCodeType(req.Type) {
			response.Response(r, w, nil, errors.New("验证码类型不正确"))
			return
		}

		// 检查发送频率限制（60秒内只能发送一次）
		rateLimitKey := fmt.Sprintf("phone_rate_limit_%s", req.Phone)
		exists, err := svcCtx.Redis.Exists(rateLimitKey).Result()
		if err != nil {
			response.Response(r, w, nil, errors.New("服务内部异常"))
			return
		}
		if exists == 1 {
			response.Response(r, w, nil, errors.New("发送过于频繁，请60秒后再试"))
			return
		}

		l := logic.NewGetPhoneCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetPhoneCode(&req)
		response.Response(r, w, resp, err)
	}
}
