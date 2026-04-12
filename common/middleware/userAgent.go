package middleware

import (
	"beaver/utils/device"
	"context"
	"net/http"
)

func UserAgentMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ua := r.Header.Get("User-Agent")
		preciseType := device.GetDeviceType(ua)
		deviceGroup := device.GetDeviceGroup(preciseType)

		// 将原始 UA、精准类型、族群全部注入 Context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user-agent", ua)
		ctx = context.WithValue(ctx, "precise-type", preciseType)
		ctx = context.WithValue(ctx, "device-group", deviceGroup)

		next(w, r.WithContext(ctx))
	}
}
