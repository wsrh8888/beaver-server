package redis_ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// WSConnectionInfo WebSocket连接信息
type WSConnectionInfo struct {
	/**
	 *  userId: 用户ID
	 */
	UserID string `json:"userId"`
	/**
	 *  deviceType: 设备类型
	 */
	DeviceType string `json:"deviceType"`
	/**
	 *  serverId: 服务器ID
	 */
	ServerID string `json:"serverId"`
	/**
	 *  connId: 连接ID
	 */
	ConnID string `json:"connId"`
	/**
	 *  loginTime: 登录时间
	 */
	LoginTime int64 `json:"loginTime"`
}

// WSManager Redis WebSocket管理器
type WSManager struct {
	client   *redis.Client
	serverID string
}

// NewWSManager 创建Redis WebSocket管理器
func NewWSManager(client *redis.Client, serverID string) *WSManager {
	return &WSManager{
		client:   client,
		serverID: serverID,
	}
}

// AddConnection 添加连接
func (m *WSManager) AddConnection(userID, deviceType, connID string) error {
	key := fmt.Sprintf("ws:conn:%s:%s", userID, deviceType)
	info := WSConnectionInfo{
		UserID:     userID,
		DeviceType: deviceType,
		ServerID:   m.serverID,
		ConnID:     connID,
		LoginTime:  time.Now().Unix(),
	}

	data, _ := json.Marshal(info)
	return m.client.HSet(key, connID, data).Err()
}

// RemoveConnection 移除连接
func (m *WSManager) RemoveConnection(userID, deviceType, connID string) error {
	key := fmt.Sprintf("ws:conn:%s:%s", userID, deviceType)
	return m.client.HDel(key, connID).Err()
}

// GetUserConnections 获取用户连接
func (m *WSManager) GetUserConnections(userID string) ([]WSConnectionInfo, error) {
	pattern := fmt.Sprintf("ws:conn:%s:*", userID)
	keys, err := m.client.Keys(pattern).Result()
	if err != nil {
		return nil, err
	}

	var connections []WSConnectionInfo
	for _, key := range keys {
		conns, err := m.client.HGetAll(key).Result()
		if err != nil {
			continue
		}

		for _, data := range conns {
			var info WSConnectionInfo
			if json.Unmarshal([]byte(data), &info) == nil {
				connections = append(connections, info)
			}
		}
	}

	return connections, nil
}

// SendToUser 发送消息给用户（跨服务器）
func (m *WSManager) SendToUser(userID string, message []byte) error {
	connections, err := m.GetUserConnections(userID)
	if err != nil {
		return err
	}

	for _, conn := range connections {
		if conn.ServerID != m.serverID {
			// 发送到其他服务器的消息队列
			streamKey := fmt.Sprintf("ws:stream:%s", conn.ServerID)
			m.client.XAdd(&redis.XAddArgs{
				Stream: streamKey,
				Values: map[string]interface{}{
					"connID":  conn.ConnID,
					"message": string(message),
				},
			})
		}
	}

	return nil
}
