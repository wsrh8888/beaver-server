package ymsgw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"

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
	token := q.Get("token")
	if token == "" {
		return ygws.ConnectionConfig{}, false
	}
	claims, err := a.verifyTokenClaims(token, q.Get("platform"))
	if err != nil {
		return ygws.ConnectionConfig{}, false
	}

	// 房间名由服务端按 docId@epoch 改写进 URL path；query.docId 是权威业务 ID。
	// path.Base 兼容调用方未改写 path 的情况（可能仍带 /api/doc/v1/ws 前缀）。
	roomDocID, roomEpoch := ParseRoom(path.Base(strings.TrimSuffix(r.URL.Path, "/")))
	docID := q.Get("docId")
	if docID == "" {
		docID = roomDocID
	}
	if docID == "" || (roomDocID != "" && roomDocID != docID) {
		return ygws.ConnectionConfig{}, false
	}

	// userId 允许从 JWT 兜底：旧客户端 WS query 未带 userId 时不至于必然 401，
	// 但仍强制要求与 token 内的 userId 一致，不能冒用他人身份。
	userID := q.Get("userId")
	if userID == "" {
		userID = claims.UserID
	}
	if userID == "" || claims.UserID != userID {
		return ygws.ConnectionConfig{}, false
	}

	// 版本回滚会递增 epoch；旧房间一律拒绝，避免旧 CRDT 状态被重新写回。
	if roomEpoch > 0 {
		ep, eerr := a.doc.ResolveRoomEpoch(context.Background(), &doc_rpc.ResolveRoomEpochReq{
			DocId: docID,
			Epoch: roomEpoch,
		})
		if eerr != nil || !ep.GetValid() {
			return ygws.ConnectionConfig{}, false
		}
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

// verifyTokenClaims 复刻 ws_api 的 WS 鉴权：JWT 签名校验 + Redis 登录态比对 +
// 设备指纹一致性。返回解析出的 claims，供上层把 userId 与 query 参数/token 对齐。
func (a *wsAuth) verifyTokenClaims(token, platform string) (*jwts.CustomClaims, error) {
	if token == "" {
		return nil, errors.New("token不能为空")
	}
	claims, err := jwts.ParseToken(token, a.secret)
	if err != nil {
		return nil, errors.New("无效的token")
	}
	platforms := []string{platform}
	if platform == "" {
		platforms = []string{"desktop", "mobile"}
	}
	for _, dt := range platforms {
		key := fmt.Sprintf("user_authentication_session:%s:%s", claims.UserID, dt)
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
					return nil, errors.New("设备标识符不匹配，疑似冒用")
				}
			}
			return claims, nil
		}
	}
	return nil, errors.New("登录已过期或在其他设备登录")
}
