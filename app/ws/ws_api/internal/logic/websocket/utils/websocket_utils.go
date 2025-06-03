package websocket_utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	ws_response "beaver/app/ws/ws_api/response"
	type_struct "beaver/app/ws/ws_api/types"
	"beaver/common/wsEnum/wsCommandConst"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var UserOnlineWsMap = make(map[string]*UserWsInfo)
var WsMapMutex sync.RWMutex // 导出互斥锁以供其他包使用

type UserWsInfo struct {
	WsClientMap map[string]*websocket.Conn // 用户管理的所有 WebSocket 通道
}

// GetUserKey 生成用户连接的唯一key
func GetUserKey(userID string, deviceType string) string {
	return userID + "_" + deviceType
}

// SendMsgToReceiverAndSyncToSender 发消息给接收者，并同步给发送者的其他设备
func SendMsgToReceiverAndSyncToSender(
	revUserID string,
	sendUserID string,
	command wsCommandConst.Command,
	receiverContent type_struct.WsContent,
	senderSyncContent type_struct.WsContent,
	excludeConn *websocket.Conn, // 排除发送消息的连接，避免重复
) {
	logx.Infof("🚀🚀🚀 调用新版本SendMsgToReceiverAndSyncToSender函数 发送者: %s 接收者: %s 🚀🚀🚀", sendUserID, revUserID)

	// 遍历所有设备类型
	deviceTypes := []string{"mobile", "windows", "mac", "linux", "web"}

	// 调试：打印发送者的所有连接状态
	logx.Infof("=== 消息发送前 发送者连接状态 用户ID: %s ===", sendUserID)
	WsMapMutex.RLock()
	for _, dt := range deviceTypes {
		userKey := GetUserKey(sendUserID, dt)
		if sendUser, ok := UserOnlineWsMap[userKey]; ok {
			logx.Infof("发送者设备类型: %s, userKey: %s, 连接数: %d", dt, userKey, len(sendUser.WsClientMap))
		}
	}
	WsMapMutex.RUnlock()
	logx.Infof("=== 发送者连接状态结束 ===")

	// 发送给接收者B的所有设备
	WsMapMutex.RLock()
	for _, deviceType := range deviceTypes {
		userKey := GetUserKey(revUserID, deviceType)
		revUser, ok := UserOnlineWsMap[userKey]
		if ok {
			jsonContent, _ := json.Marshal(receiverContent)
			logx.Info("发送消息给接收者：", revUserID, "设备类型：", deviceType, "发送者：", sendUserID, "消息内容：", string(jsonContent))
			sendWsMapMsg(revUser.WsClientMap, command, receiverContent)
		}
	}
	WsMapMutex.RUnlock()

	// 发送给发送者A的其他设备（用于多端同步）
	WsMapMutex.RLock()
	excludeAddr := ""
	if excludeConn != nil {
		excludeAddr = excludeConn.RemoteAddr().String()
	}

	for _, deviceType := range deviceTypes {
		userKey := GetUserKey(sendUserID, deviceType)
		sendUser, ok := UserOnlineWsMap[userKey]
		if ok {
			// 如果指定了要排除的连接，需要检查这个连接是否在当前设备类型中
			if excludeConn != nil {
				filteredMap := make(map[string]*websocket.Conn)
				hasExcludedConn := false

				logx.Infof("检查设备类型 %s 的连接，排除地址: %s", deviceType, excludeAddr)
				for addr, conn := range sendUser.WsClientMap {
					if addr == excludeAddr {
						hasExcludedConn = true
						logx.Infof("在设备类型 %s 中跳过发送方连接: %s", deviceType, addr)
					} else {
						filteredMap[addr] = conn
						logx.Infof("在设备类型 %s 中保留连接: %s", deviceType, addr)
					}
				}

				if hasExcludedConn && len(filteredMap) == 0 {
					logx.Infof("设备类型 %s 只有发送方连接，无需同步", deviceType)
				} else if len(filteredMap) > 0 {
					jsonContent, _ := json.Marshal(senderSyncContent)
					logx.Info("同步消息给发送者的其他设备：", sendUserID, "设备类型：", deviceType, "接收者：", revUserID, "消息内容：", string(jsonContent))
					sendWsMapMsg(filteredMap, command, senderSyncContent)
				} else {
					// 当前设备类型没有排除的连接，说明这不是发送方设备类型，全部发送
					jsonContent, _ := json.Marshal(senderSyncContent)
					logx.Info("同步消息给发送者的其他设备类型：", sendUserID, "设备类型：", deviceType, "接收者：", revUserID, "消息内容：", string(jsonContent))
					sendWsMapMsg(sendUser.WsClientMap, command, senderSyncContent)
				}
			} else {
				jsonContent, _ := json.Marshal(senderSyncContent)
				logx.Info("同步消息给发送者：", sendUserID, "设备类型：", deviceType, "接收者：", revUserID, "消息内容：", string(jsonContent))
				sendWsMapMsg(sendUser.WsClientMap, command, senderSyncContent)
			}
		} else {
			logx.Infof("发送者设备类型 %s 未找到在线连接", deviceType)
		}
	}
	WsMapMutex.RUnlock()
}

// sendWsMapMsg 给一组的 WebSocket 通道发送消息
func sendWsMapMsg(wsMap map[string]*websocket.Conn, command wsCommandConst.Command, content type_struct.WsContent) {
	for _, conn := range wsMap {
		ws_response.WsResponse(conn, command, content)
	}
}

func IsGroupChat(conversationID string) bool {
	return !strings.Contains(conversationID, "_")
}

func GetRecipientIdFromConversationID(conversationID string, userID string) string {
	ids := strings.Split(conversationID, "_")
	if len(ids) != 2 {
		fmt.Println("无效的会话Id：", conversationID)
		return ""
	}
	if ids[0] == userID {
		return ids[1]
	}
	return ids[0]
}

// GetConnectionStats 获取连接统计信息
func GetConnectionStats() map[string]interface{} {
	stats := make(map[string]interface{})

	totalUsers := len(UserOnlineWsMap)
	totalConnections := 0
	deviceStats := make(map[string]int)

	for userKey, userWsInfo := range UserOnlineWsMap {
		totalConnections += len(userWsInfo.WsClientMap)

		// 统计设备类型
		parts := strings.Split(userKey, "_")
		if len(parts) >= 2 {
			deviceType := parts[1]
			deviceStats[deviceType]++
		}
	}

	stats["total_users"] = totalUsers
	stats["total_connections"] = totalConnections
	stats["device_stats"] = deviceStats

	return stats
}

// PrintConnectionStats 打印连接统计信息
func PrintConnectionStats() {
	stats := GetConnectionStats()
	logx.Infof("=== WebSocket连接统计 ===")
	logx.Infof("在线用户数: %d", stats["total_users"])
	logx.Infof("总连接数: %d", stats["total_connections"])
	logx.Infof("设备分布: %+v", stats["device_stats"])
	logx.Infof("========================")
}

// SendMsgToUser 只发送消息给指定用户的所有设备
func SendMsgToUser(targetUserID string, command wsCommandConst.Command, content type_struct.WsContent) {
	// 遍历所有设备类型
	deviceTypes := []string{"mobile", "windows", "mac", "linux", "web"}

	WsMapMutex.RLock()
	defer WsMapMutex.RUnlock()

	for _, deviceType := range deviceTypes {
		userKey := GetUserKey(targetUserID, deviceType)
		userInfo, ok := UserOnlineWsMap[userKey]
		if ok {
			jsonContent, _ := json.Marshal(content)
			logx.Infof("发送消息给用户：%s, 设备类型：%s, 消息内容：%s", targetUserID, deviceType, string(jsonContent))
			sendWsMapMsg(userInfo.WsClientMap, command, content)
		}
	}
}
