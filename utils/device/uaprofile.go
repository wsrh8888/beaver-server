package device

import (
	"net/url"
	"strings"
)

// UAProfile 从客户端 UA 解析的设备档案
type UAProfile struct {
	PreciseType string // windows / ios / android / macos / linux
	DeviceGroup string // desktop / mobile
	Model       string // iPhone17,3、SM-G991B、Windows-PC
	OsVersion   string // 18.2、10.0.19045
	DisplayName string // iPhone 17 Pro、DESKTOP-HOME
}

func ParseUAProfile(userAgent string) UAProfile {
	preciseType := GetDeviceType(userAgent)
	return UAProfile{
		PreciseType: preciseType,
		DeviceGroup: GetDeviceGroup(preciseType),
		Model:       parseUAToken(userAgent, "model"),
		OsVersion:   parseUAToken(userAgent, "os"),
		DisplayName: parseUAToken(userAgent, "name"),
	}
}

func parseUAToken(userAgent, key string) string {
	prefix := key + "/"
	idx := strings.Index(userAgent, prefix)
	if idx < 0 {
		return ""
	}
	rest := userAgent[idx+len(prefix):]
	end := strings.IndexAny(rest, " )")
	raw := rest
	if end >= 0 {
		raw = rest[:end]
	}
	decoded, _ := url.QueryUnescape(strings.TrimSpace(raw))
	return decoded
}
