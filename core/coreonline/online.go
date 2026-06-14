// Package coreonline 维护 WS 实时在线态（Redis，带 TTL 的 ephemeral 数据）。
//
//   - auth_models.AuthDeviceModel（MySQL）：设备档案 —— 持久化
//   - coreonline（Redis）：当前 WS 是否连着 —— 实时、可过期
//   - ws_api UserOnlineWsMap：本机连接，本实例推送
//
// MarkOnline / MarkOffline 由 ws_api 在连/断/心跳时内部调用，不是 HTTP 接口。
package coreonline

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	keyPrefix = "beaver:user:online:"
	// 须大于 ws.yaml WebSocket.PingPeriod（240s），心跳用 MarkOnline 续期。
	ttl = 6 * time.Minute
)

type Slot struct {
	InstanceID string
	Slot       string // desktop | mobile
}

type User struct {
	UserID string
	Slots  []Slot
}

func key(userID string) string {
	return keyPrefix + userID
}

func parseMember(member string) (instanceID, slot string) {
	i := strings.LastIndex(member, ":")
	if i <= 0 {
		return member, ""
	}
	return member[:i], member[i+1:]
}

// MarkOnline 连接建立或心跳续期（SADD + EXPIRE）。
func MarkOnline(rdb *redis.Client, userID, slot, instanceID string) {
	if rdb == nil || userID == "" || slot == "" || instanceID == "" {
		return
	}
	k := key(userID)
	rdb.SAdd(k, instanceID+":"+slot)
	rdb.Expire(k, ttl)
}

// MarkOffline 连接断开。
func MarkOffline(rdb *redis.Client, userID, slot, instanceID string) {
	if rdb == nil || userID == "" || slot == "" || instanceID == "" {
		return
	}
	k := key(userID)
	rdb.SRem(k, instanceID+":"+slot)
	if rdb.SCard(k).Val() == 0 {
		rdb.Del(k)
	}
}

// IsSlotOnline 某槽位（desktop/mobile）是否 WS 在线。
func IsSlotOnline(rdb *redis.Client, userID, slot string) bool {
	if rdb == nil || userID == "" || slot == "" {
		return false
	}
	members, err := rdb.SMembers(key(userID)).Result()
	if err != nil {
		return false
	}
	suffix := ":" + slot
	for _, m := range members {
		if strings.HasSuffix(m, suffix) {
			return true
		}
	}
	return false
}

// IsOnline 是否在线。
func IsOnline(rdb *redis.Client, userID string) bool {
	if rdb == nil || userID == "" {
		return false
	}
	return rdb.SCard(key(userID)).Val() > 0
}

// List 扫描全部在线用户（后台监控）。
func List(rdb *redis.Client) ([]User, error) {
	if rdb == nil {
		return nil, nil
	}
	list := make([]User, 0)
	var cursor uint64
	for {
		keys, next, err := rdb.Scan(cursor, keyPrefix+"*", 200).Result()
		if err != nil {
			return list, err
		}
		for _, k := range keys {
			members, err := rdb.SMembers(k).Result()
			if err != nil {
				return list, err
			}
			u := User{UserID: strings.TrimPrefix(k, keyPrefix)}
			for _, m := range members {
				inst, slot := parseMember(m)
				u.Slots = append(u.Slots, Slot{InstanceID: inst, Slot: slot})
			}
			list = append(list, u)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return list, nil
}
