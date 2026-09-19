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

package ws

import (
	"context"
	"net/http"
	"time"

	ws "beaver/app/open/open_api/internal/logic/ws/websocket"
	ws_auth "beaver/app/open/open_api/internal/logic/ws/websocket/auth"
	ws_conn "beaver/app/open/open_api/internal/logic/ws/websocket/conn"
	"beaver/app/open/open_api/internal/logic/ws/websocket/heartbeat"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/gorilla/websocket"
)

// OpenWsLogic 开放平台机器人长连接入口。
//
// 连接生命周期：鉴权（appId + ticket）→ 升级 → 注册到独立注册表 → 心跳 + 消息循环。
// 与用户侧 ws_api 的关键区别：不调用 coreonline，机器人不进入「在线用户」集合。
type OpenWsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

// 建立长连接（WebSocket 升级；本端点不返回 JSON，升级后走 WS 帧协议）
func NewOpenWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OpenWsLogic {
	return &OpenWsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("open_ws", ctx),
	}
}

func (l *OpenWsLogic) OpenWs(req *types.OpenWsReq, w http.ResponseWriter, r *http.Request) (resp *types.OpenWsRes, err error) {
	// 1. 鉴权先于升级执行：校验失败直接以 HTTP 状态码拒绝，不建立连接
	identity, authErr := ws_auth.VerifyRobotToken(l.svcCtx.DB, req.AppID, req.Ticket)
	if authErr != nil {
		l.logger.Error(model.LogMsg{
			Text: "机器人长连接鉴权失败",
			Data: map[string]any{"appId": req.AppID, "remoteAddr": r.RemoteAddr, "err": authErr.Error()},
		})
		http.Error(w, authErr.Error(), http.StatusUnauthorized)
		return nil, nil
	}

	// 2. 升级 HTTP → WebSocket
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	wsConn, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "机器人长连接升级失败",
			Data: map[string]any{"appId": identity.AppID, "err": err.Error()},
		})
		return nil, nil
	}

	// 3. 配置连接参数（读超时 + Pong 续期）
	l.configureConn(wsConn)

	client := ws_conn.NewClient(wsConn)
	addr := wsConn.RemoteAddr().String()

	// 4. 注册到独立注册表（key = robotId），同一机器人重复建连会踢掉旧连接
	if kicked := ws_conn.Register(identity.RobotID, client); kicked > 0 {
		l.logger.Info(model.LogMsg{
			Text: "机器人长连接已接管旧连接",
			Data: map[string]any{"robotId": identity.RobotID, "kicked": kicked},
		})
	}
	l.logger.Info(model.LogMsg{
		Text: "机器人长连接已建立",
		Data: map[string]any{
			"appId":      identity.AppID,
			"robotId":    identity.RobotID,
			"remoteAddr": addr,
		},
	})

	defer func() {
		_ = wsConn.Close()
		ws_conn.Unregister(identity.RobotID, addr)
		l.logger.Info(model.LogMsg{
			Text: "机器人长连接已断开",
			Data: map[string]any{"appId": identity.AppID, "robotId": identity.RobotID, "remoteAddr": addr},
		})
	}()

	// 5. 启动服务端主动心跳
	hbManager := heartbeat.NewManager(client, identity.RobotID, l.pingPeriod(), l.writeWait())
	defer hbManager.Stop()
	hbManager.Start()

	// 6. 消息循环（阻塞至连接关闭）
	ws.HandleRobotMessages(l.ctx, l.svcCtx, identity.RobotID, client)

	return nil, nil
}

// configureConn 设置读限制与读超时；收到 Pong 帧即续期。
func (l *OpenWsLogic) configureConn(wsConn *websocket.Conn) {
	cfg := l.svcCtx.Config.WebSocket
	if cfg.MaxMessageSize > 0 {
		wsConn.SetReadLimit(int64(cfg.MaxMessageSize))
	}
	pongWait := l.pongWait()
	_ = wsConn.SetReadDeadline(time.Now().Add(pongWait))
	wsConn.SetPongHandler(func(string) error {
		return wsConn.SetReadDeadline(time.Now().Add(pongWait))
	})
}

func (l *OpenWsLogic) pongWait() time.Duration {
	if v := l.svcCtx.Config.WebSocket.PongWait; v > 0 {
		return time.Duration(v) * time.Second
	}
	return 300 * time.Second
}

func (l *OpenWsLogic) pingPeriod() time.Duration {
	if v := l.svcCtx.Config.WebSocket.PingPeriod; v > 0 {
		return time.Duration(v) * time.Second
	}
	return 30 * time.Second
}

func (l *OpenWsLogic) writeWait() time.Duration {
	if v := l.svcCtx.Config.WebSocket.WriteWait; v > 0 {
		return time.Duration(v) * time.Second
	}
	return 10 * time.Second
}
