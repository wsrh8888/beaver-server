package handler

import (
	"errors"

	logic "beaver/app/auth/auth_api/internal/logic/auth_public"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/middleware/ua"
	"beaver/common/response"
	"beaver/utils/device"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PhoneLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PhoneLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}
		if err := validateLoginDevice(r, req.DeviceID); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		l := logic.NewPhoneLoginLogic(r.Context(), svcCtx)
		resp, err := l.PhoneLogin(&req)
		response.Response(r, w, resp, err)
	}
}

func validateLoginDevice(r *http.Request, deviceID string) error {
	preciseType := ua.DeviceType(r.Context())
	if preciseType == "" || preciseType == device.DeviceUnknown {
		return errors.New("不支持的设备类型")
	}
	if deviceID == "" {
		return errors.New("无法识别的物理设备，请联系管理员")
	}
	return nil
}
