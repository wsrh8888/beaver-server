package auth_public

import (
	"context"

	"beaver/common/middleware/ua"
	"beaver/utils/device"
)

func ctxUAProfile(ctx context.Context) device.UAProfile {
	return ctx.Value(ua.KeyUAProfile).(device.UAProfile)
}

func ctxClientIP(ctx context.Context) string {
	return ctx.Value("ClientIP").(string)
}
