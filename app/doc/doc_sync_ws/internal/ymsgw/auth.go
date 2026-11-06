package ymsgw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	doccli "beaver/app/doc/doc_rpc/doc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	"beaver/utils/jwts"

	"github.com/go-redis/redis"
	ygws "github.com/reearth/ygo/provider/websocket"
)

// wsAuth 封装升级前鉴权所需依赖。
type wsAuth struct {
	secret string
	rdb    *redis.Client
	doc    doccli.Doc
}

// authorize 返回 ygo Server.Authorize 钩子：升级前完成 JWT + 登录态校验与文档权限判定。
// 校验失败返回 ok=false（ygo 会以 401 拒绝，不升级连接）；viewer/commenter 置只读。
func (a *wsAuth) authorize(r *http.Request) (ygws.ConnectionConfig, bool) {
	q := r.URL.Query()
	docID := q.Get("docId")
	userID := q.Get("userId")
	token := q.Get("token")
	platform := q.Get("platform")
	if docID == "" || userID == "" || token == "" {
		return ygws.ConnectionConfig{}, false
	}
	if err := a.verifyToken(token, userID, platform); err != nil {
		return ygws.ConnectionConfig{}, false
	}
	res, err := a.doc.CheckPermission(context.Background(), &doc_rpc.CheckPermissionReq{
		ResourceType: "doc",
		ResourceId:   docID,
		UserId:       userID,
	})
	if err != nil {
		return ygws.ConnectionConfig{}, false
	}
	role := res.GetRole()
	if roleRank(role) < roleRank("viewer") {
		// 无任何有效角色：拒绝连接。
		return ygws.ConnectionConfig{}, false
	}
	// editor 及以上可写，commenter/viewer 只读。
	return ygws.ConnectionConfig{ReadOnly: roleRank(role) < roleRank("editor")}, true
}

// verifyToken 复刻 ws_api 的 WS 鉴权：JWT 签名校验 + Redis 登录态比对 + 设备指纹一致性。
func (a *wsAuth) verifyToken(token, userID, platform string) error {
	if token == "" {
		return errors.New("token不能为空")
	}
	claims, err := jwts.ParseToken(token, a.secret)
	if err != nil {
		return errors.New("无效的token")
	}
	if claims.UserID != userID {
		return errors.New("token与用户不匹配")
	}
	platforms := []string{platform}
	if platform == "" {
		platforms = []string{"desktop", "mobile"}
	}
	for _, dt := range platforms {
		key := fmt.Sprintf("user_authentication_session:%s:%s", userID, dt)
		val, err := a.rdb.Get(key).Result()
		if err != nil || val == "" {
			continue
		}
		var info map[string]any
		if json.Unmarshal([]byte(val), &info) != nil {
			continue
		}
		if storedToken, _ := info["token"].(string); storedToken == token {
			if claims.DeviceID != "" {
				if storedDeviceID, ok := info["device_id"].(string); ok && storedDeviceID != claims.DeviceID {
					return errors.New("设备标识符不匹配，疑似冒用")
				}
			}
			return nil
		}
	}
	return errors.New("登录已过期或在其他设备登录")
}
