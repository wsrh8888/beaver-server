/*
 * Copyright (c) 2024-2026 Beaver IM Team
 * SPDX-License-Identifier: MIT
 * Project: beaver-server
 * https://github.com/wsrh8888/beaver-server
 *
 * 中文：
 * 本文件为海狸 IM（Beaver IM）开源项目源代码。
 * 版权所有 © 2024-2026 Beaver IM Team，基于 MIT 协议授权。
 * 禁止删除、篡改或替换本文件头部版权与许可声明。
 * 使用与商业授权说明：https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * English:
 * This file is part of the Beaver IM open-source project.
 * Copyright (c) 2024-2026 Beaver IM Team. Licensed under the MIT License.
 * Do not remove, alter, or replace this copyright and license header.
 * Usage & commercial licensing: https://wsrh8888.github.io/beaver-docs/community/license.html
 *
 * beaver-server-header-v1
 */

package logic

import (
	"context"
	"net/http"
	"time"

	ws "beaver/app/ws/ws_api/internal/logic/websocket"
	ws_auth "beaver/app/ws/ws_api/internal/logic/websocket/auth"
	ws_conn "beaver/app/ws/ws_api/internal/logic/websocket/conn"
	"beaver/app/ws/ws_api/internal/logic/websocket/heartbeat"
	"beaver/app/ws/ws_api/internal/svc"
	"beaver/app/ws/ws_api/internal/types"
	"beaver/core/coreonline"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
	"beaver/utils/device"

	"github.com/gorilla/websocket"
)

type ChatWebsocketLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewChatWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChatWebsocketLogic {
	return &ChatWebsocketLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("chat_websocket", ctx),
	}
}
func (l *ChatWebsocketLogic) ChatWebsocket(req *types.WsReq, w http.ResponseWriter, r *http.Request) (resp *types.WsRes, err error) {
	// 1. 基础入口校验：精准识别设备
	userAgent := r.Header.Get("User-Agent")
	preciseType := device.GetDeviceType(userAgent)
	if preciseType == device.DeviceUnknown {
		l.logger.Error(model.LogMsg{
			Text: "连接拒绝非法设备接入",
			Data: map[string]any{"userId": req.UserID, "userAgent": userAgent},
		})
		http.Error(w, "Illegal Device", http.StatusForbidden)
		return nil, nil
	}

	// 2. 获取所属槽位 (Group: mobile/desktop)，用于检索 Redis 登录态
	deviceGroup := device.GetDeviceGroup(preciseType)

	// 3. 鉴权：先于升级执行。使用槽位 (Group) 进行登录态比对
	if authErr := ws_auth.VerifyWsToken(req.Token, l.svcCtx.Config.Auth.AccessSecret, req.UserID, deviceGroup, l.svcCtx.Redis); authErr != nil {
		l.logger.Error(model.LogMsg{
			Text: "WS鉴权失败",
			Data: map[string]any{"userId": req.UserID, "preciseType": preciseType, "deviceGroup": deviceGroup, "err": authErr.Error()},
		})
		http.Error(w, authErr.Error(), http.StatusUnauthorized)
		return nil, nil
	}

	// 2. 升级 HTTP → WebSocket
	conn, err := upgradeToWebSocket(w, r)
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "WebSocket升级失败",
			Data: map[string]any{"userId": req.UserID, "err": err.Error()},
		})
		return nil, nil
	}

	// 3. 配置连接参数
	configureWebSocketConn(conn, l.svcCtx)

	// 4. 封装为 Client（带写 mutex）
	client := ws_conn.NewClient(conn)

	// 5. 注册连接：统一使用槽位 (Group) 管理，确保互踢逻辑闭合
	// 无论具体的 OS 是 windows、macos 还是 linux，在 WS 路由层统一视为 desktop 槽位
	userKey := ws_conn.GetUserKey(req.UserID, deviceGroup)

	l.logger.Info(model.LogMsg{
		Text: "用户上线",
		Data: map[string]interface{}{
			"userId":      req.UserID,
			"deviceGroup": deviceGroup,
			"preciseType": preciseType,
			"remoteAddr":  conn.RemoteAddr().String(),
		},
	})
	l.manageUserConnection(userKey, client, req.UserID, deviceGroup)
	coreonline.MarkOnline(l.svcCtx.Redis, req.UserID, deviceGroup, l.svcCtx.InstanceID)
	connAddr := conn.RemoteAddr().String()
	defer func() {
		conn.Close()
		l.cleanupConnection(req.UserID, deviceGroup, connAddr)
	}()

	// 6. 启动心跳
	heartbeatManager := heartbeat.NewManager(client, req.UserID, deviceGroup, l.svcCtx)
	defer heartbeatManager.Stop()
	heartbeatManager.Start()

	// 7. 消息循环
	ws.HandleWebSocketMessages(l.ctx, l.svcCtx, req, deviceGroup, r, client)

	return nil, nil
}

func (l *ChatWebsocketLogic) manageUserConnection(userKey string, client *ws_conn.Client, userID, deviceGroup string) {
	ws_conn.WsMapMutex.Lock()
	defer ws_conn.WsMapMutex.Unlock()

	addr := client.Conn.RemoteAddr().String()
	userWsInfo, ok := ws_conn.UserOnlineWsMap[userKey]

	if ok {
		// 槽位级互踢：目前 desktop 和 mobile 均限制单物理设备在线
		if deviceGroup == "desktop" || deviceGroup == "mobile" {
			for oldAddr, oldClient := range userWsInfo.WsClientMap {
				l.logger.Info(model.LogMsg{
					Text: "槽位互踢关闭旧连接",
					Data: map[string]interface{}{"userId": userID, "deviceGroup": deviceGroup, "oldAddr": oldAddr},
				})
				oldClient.Conn.Close()
				delete(userWsInfo.WsClientMap, oldAddr)
			}
		}
		userWsInfo.WsClientMap[addr] = client
	} else {
		ws_conn.UserOnlineWsMap[userKey] = &ws_conn.UserWsInfo{
			WsClientMap: map[string]*ws_conn.Client{addr: client},
		}
	}

	l.logger.Info(model.LogMsg{
		Text: "连接注册成功",
		Data: map[string]interface{}{"userId": userID, "deviceGroup": deviceGroup},
	})
}

func configureWebSocketConn(conn *websocket.Conn, svcCtx *svc.ServiceContext) {
	conn.SetReadLimit(int64(svcCtx.Config.WebSocket.MaxMessageSize))
	pongWait := time.Duration(svcCtx.Config.WebSocket.PongWait) * time.Second
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
}

func upgradeToWebSocket(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return upGrader.Upgrade(w, r, nil)
}

func (l *ChatWebsocketLogic) cleanupConnection(userID, deviceGroup, addr string) {
	ws_conn.WsMapMutex.Lock()
	defer ws_conn.WsMapMutex.Unlock()

	userKey := ws_conn.GetUserKey(userID, deviceGroup)
	userWsInfo, ok := ws_conn.UserOnlineWsMap[userKey]
	if !ok {
		coreonline.MarkOffline(l.svcCtx.Redis, userID, deviceGroup, l.svcCtx.InstanceID)
		return
	}

	delete(userWsInfo.WsClientMap, addr)
	if len(userWsInfo.WsClientMap) == 0 {
		delete(ws_conn.UserOnlineWsMap, userKey)
		coreonline.MarkOffline(l.svcCtx.Redis, userID, deviceGroup, l.svcCtx.InstanceID)
		l.logger.Info(model.LogMsg{
			Text: "用户下线",
			Data: map[string]interface{}{"userId": userID, "deviceGroup": deviceGroup},
		})
	}
}
