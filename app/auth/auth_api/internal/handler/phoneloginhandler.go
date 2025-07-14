package handler

import (
	"beaver/app/auth/auth_api/internal/logic"
	"beaver/app/auth/auth_api/internal/svc"
	"beaver/app/auth/auth_api/internal/types"
	"beaver/common/response"
	"beaver/common/validator"
	"beaver/utils/device"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 获取登录失败次数
func getLoginFailCount(svcCtx *svc.ServiceContext, phone string) (int, error) {
	key := fmt.Sprintf("login_fail_%s", phone)
	count, err := svcCtx.Redis.Get(key).Int()
	if err != nil {
		return 0, nil
	}
	return count, nil
}

// 增加登录失败次数
func incLoginFailCount(svcCtx *svc.ServiceContext, phone string) error {
	key := fmt.Sprintf("login_fail_%s", phone)
	return svcCtx.Redis.Incr(key).Err()
}

// 重置登录失败次数
func resetLoginFailCount(svcCtx *svc.ServiceContext, phone string) error {
	key := fmt.Sprintf("login_fail_%s", phone)
	return svcCtx.Redis.Del(key).Err()
}

func phoneLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PhoneLoginReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 验证登录参数
		if err := validator.ValidateLoginParams(req.Phone, req.Password); err != nil {
			response.Response(r, w, nil, err)
			return
		}

		// 验证设备ID
		if !device.IsValidDeviceID(req.DeviceID) {
			response.Response(r, w, nil, errors.New("设备ID格式不正确"))
			return
		}

		// 检查登录失败次数
		failCount, err := getLoginFailCount(svcCtx, req.Phone)
		if err != nil {
			response.Response(r, w, nil, errors.New("服务内部异常"))
			return
		}
		if failCount >= 10 {
			response.Response(r, w, nil, errors.New("登录失败次数过多,请稍后再试"))
			return
		}

		l := logic.NewPhoneLoginLogic(r.Context(), svcCtx)
		resp, err := l.PhoneLogin(&req)
		if err != nil {
			// 增加登录失败次数
			if err := incLoginFailCount(svcCtx, req.Phone); err != nil {
			}
			response.Response(r, w, nil, err)
			return
		}

		// 登录成功，重置失败次数
		if err := resetLoginFailCount(svcCtx, req.Phone); err != nil {
			// 重置失败次数失败不影响登录成功
			logx.Errorf("重置登录失败次数失败: %v", err)
		}

		response.Response(r, w, resp, nil)
	}
}
