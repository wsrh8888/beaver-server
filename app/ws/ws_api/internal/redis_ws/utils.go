package redis_ws

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

// GetWSManager 获取WebSocket管理器实例
func GetWSManager(client *redis.Client) *WSManager {
	serverID := fmt.Sprintf("%s_%d", getHostname(), os.Getpid())
	return NewWSManager(client, serverID)
}

// getHostname 获取主机名
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
