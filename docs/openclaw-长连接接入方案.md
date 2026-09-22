# 海狸 IM 接入 OpenClaw（龙虾）长连接方案

> 目标：让 OpenClaw 作为「渠道」，通过**长连接**接入海狸 IM。
> 用户在 IM 里发消息 → 推到 OpenClaw → OpenClaw 回复 → 回到 IM。
> 机器人主动出站连接，不依赖公网 IP、域名、TLS 证书、内网穿透。
>
> **状态：服务端已实现（M1 + M2 已打通），插件侧已同步。**

---

## 一、结论

**在开放平台服务 `open_api` 内新增一个 WebSocket 端点承载机器人长连接，与用户侧的 `ws_api` 完全分离。**

| 项 | 值 |
|---|---|
| 服务 | `open_api`（已存在的开放平台服务） |
| 入口 | `GET /api/open/ws/v1/connect?appId=..&ticket=..` |
| 连接注册表 | `BotOnlineWsMap`，按 `robotId` 索引，**不调用 `coreonline`** |
| 事件通道 | 新 MQ topic `ws_bot_push_topic` / group `open_api_bot_consumer_group` |
| 网关 | **零改动**（`/api/open/*` 整体放行） |

不复用 `ws_api`（有正确性问题），也不另起独立服务（没必要）—— 理由见第二节。

---

## 二、为什么不复用 ws_api，也不另起服务

### 1. 不能混进 `ws_api`（正确性问题，不只是"感觉混"）

`core/coreonline/online.go` 用 Redis key `beaver:user:online:<userID>` 维护在线集合，被这些地方消费：

| 消费方 | 用途 |
|---|---|
| `backend_admin/.../getonlinestatslogic.go` | 后台在线用户统计 |
| `backend_admin/.../getonlineuserlistlogic.go` | 后台在线用户列表 |
| `backend_admin/.../getdashboardoverviewlogic.go` | 仪表盘总览 |
| `chat_rpc/push_notify.go` `coreonline.IsOnline(...)` | 判断接收者是否在线，决定走 WS 还是离线推送 |
| `auth_api/.../getdeviceslogic.go` | 用户查看自己的登录设备 |

而 `ws_api/internal/logic/websocket/heartbeat/manager.go` 的心跳里会调 `coreonline.MarkOnline(...)`。
**机器人若走这条链路，robotId 会被写进「在线用户」集合** → 后台在线统计/仪表盘混入机器人，`push_notify` 也会误判。

另外 `ws_api/internal/logic/websocket/conn/manager.go` 的 `SendMsgToUser` 是遍历整个 `UserOnlineWsMap` 做前缀匹配的 O(N) 扫描，机器人混入只会让 N 更大。

### 2. 但也不必另起独立服务

早期方案曾计划新建独立服务 `open_bot_api`（路径 `/api/open_bot/v1/ws`），**已放弃**：

- 独立服务的收益是「隔离 + 可独立扩缩容」。其中**隔离**这一点，只要满足下面三条就已经拿到，与放在哪个进程无关：
  1. 独立连接注册表（`BotOnlineWsMap`，不复用 `UserOnlineWsMap`）
  2. 不调用 `coreonline`
  3. 独立 MQ topic（`ws_bot_push_topic`，用户侧 `ws_api` 完全看不到机器人事件）
- `open_api` **本身就不是用户连接服务**，机器人连接放进来不会碰 `coreonline`，也不与任何用户连接共享注册表 —— 当初「不能混进 `ws_api`」的判断依然成立。
- 代价对比：独立服务要多一个进程、一份 etcd 注册、一条 `PublicList` 放行、一份部署编排；放 `open_api` 则是**零新增服务 + 零网关改动**。

**唯一代价**：长连接与开放平台 HTTP API 同进程，不能独立扩缩容。当机器人连接数增长到与在线人数同量级（知音楼那种「每人本机一个」形态）时，需重新评估是否拆出独立进程。届时是**搬迁**（注册表 + 消费者 + 路由），协议不变。

---

## 三、概念对齐（先统一叫法）

「bot / robot / agent / 私聊群聊」之所以感觉乱，是因为三个词被混用。现有模型其实是两个正交维度。

### 维度一：会话类型 `conversationType` —— 传输结构

| 值 | 类型 | 会话 ID 形态 |
|---|---|---|
| 1 | 私聊 | `private_<userId1>_<userId2>` 或 `<userId1>_<userId2>` |
| 2 | 群聊 | `group_<uuid>` |
| 3 | 圈子 | `circle_<uuid>` |

来源：`utils/conversation` 的 `ParseConversationWithType` / `GetConversationType`。

### 维度二：用户类型 `UserType` —— 说话的人是谁

| 值 | 类型 | 能力 | 对应 API | 定义位置 |
|---|---|---|---|---|
| 1 | 普通用户 | — | — | — |
| 2 | **通知机器人** | **单向**：只能发，收不到回复 | `bot_public.api` | `user_models.UserTypeBot` |
| 3 | **智能机器人** | **双向**：能收事件、能回消息 | `robot.api` + 本长连接 | `user_models.UserTypeRobot` |

### 关键结论

1. **「agent 聊天」不是新的会话类型。** 它就是会话类型 1 或 2，对方是 `UserType=3` 的智能机器人。**agent 在机器人后面，不在会话结构里。**
2. **「bot 聊天」严格说不成立。** `UserType=2` 是单向推送，收不到回复。
3. 统一叫法：**通知机器人**(2) / **智能机器人**(3) / **Agent**（非 `user_type`，是产品概念）。

### 一个应用 = 一个机器人

`app/open/open_models/open_robot_model.go`：`OpenRobot.AppID` 带 `uniqueIndex`，注释写明「一个 App 对应一个 Robot IM 用户」。

- 单个应用只能有一个智能机器人（与飞书/企微一致）
- 整个部署可以有 N 个机器人（N 应用 → N 机器人），`BotOnlineWsMap` 按 `robotId` 索引，天然支持多条连接并存
- 若将来要「一应用多机器人」，需移除 `AppID` 的 `uniqueIndex`，但那会让权限/配额模型变成 1:N，**不建议**

### 一个群里可以放几个机器人？—— 任意个，且无需改代码

`chat_rpc/internal/logic/robot_webhook_push.go` 的 `pushGroupAt` 遍历 `msg.AtUserIDs`，对**每一个**被 @ 的用户查是不是机器人，是就各自投递一次。外部机器人与自有 agent 可同群，各收各的 @。

> 注意：平台**不**存储「机器人是否在单聊/群聊响应」这类开关。
> 平台只负责把消息送达，不替机器人决定该不该响应 —— 机器人不想回，自己忽略即可。
> （这曾导致一个隐藏 bug：旧的 `!res.EnableSingleChat` 判断恒为真，事件永远发不出去，已修复。）

### 核心心智模型：IM 不知道「机器人后面是谁」

| | 外部（如 OpenClaw） | 你自己的 agent |
|---|---|---|
| IM 里的身份 | `UserType=3` 的用户 | `UserType=3` 的用户 |
| 私聊 1:1 | 同一机制 | 一样 |
| 群里 @ | 同一机制 | 一样 |
| 接入接口 | 开放平台 `appId` + 长连接 | **同一套** |
| **唯一区别** | 机器人后面接的是用户的 OpenClaw | 机器人后面接的是 `beaver-agent` |

**结论：你自己的 agent 不是特殊形态，它就是「又一个接入方」。**
由此：不需要为 agent 新增会话类型、不需要特殊通道、**开源侧（IM + 开放平台）与商业侧（agent）的边界就是这条协议本身**。

若希望 agent「看起来是内置的」：**预置一个官方应用**（`appId`/`appSecret` 内置在客户端），而不是给特权通道。

---

## 四、网关为什么零改动

网关**按路径段做服务发现**：`app/gateway/gateway_api/core/proxy.go` 用正则 `/api/(.*?)/` 取段名，拼 `_api` 去 etcd 查地址。

- `/api/open/ws/v1/connect` → 段名 `open` → 查 `open_api` → **自动路由，无需配置路由表**
- `p.auth()` 的五步顺序里，`isOpenApiPassThrough` 对 `/api/open/*` **整体放行** → 机器人没有用户 JWT 也不会被 401

因此**不需要**动 `gateway.yaml` 的 `PublicList`。

> 对比：若当初选了 `/api/open_bot/*` 或 `/api/open_ws/*`（不属于 `/api/open/*` 前缀），则**必须**加一行 `PublicList` 放行。这也是选 `open_api` 内实现的一个实际收益。

---

## 五、目标架构

```
IM 用户发消息
   │
   ▼
chat_rpc.sendmsglogic ──► robot_webhook_push.tryPush()
                                  │
                          dispatch(robotId, appId, eventType, event)
                                  │
                    投 MqTopicWsBot（targetId = robotId）
                                  │
                                  ▼
        open_api 内的消费者（broadcast 模式，每个实例都消费）
                                  │
                   BotOnlineWsMap（按 robotId 索引，不写 coreonline）
                                  │
                        本实例没有该机器人的连接？→ 丢弃
                                  │ 有
                                  ▼
                  机器人 wss 连接（客户端主动出站）
                                  │
                                  ▼
                        OpenClaw 插件 → AI → 回复
                                  │
        POST /api/open/robot/v1/send_message（已有接口，未改动）
                                  ▼
                            回到 IM 会话
```

**回复那条腿本来就是通的**，本次改的只是「事件下行」这一条腿。

---

## 六、实现清单（已完成）

### 服务端（beaver-server）

| 文件 | 说明 |
|---|---|
| `common/const/mqwsconst/mqWsConst.go` | 新增 `MqTopicWsBot = "ws_bot_push_topic"`、`MqGroupWsBot = "open_api_bot_consumer_group"` |
| `app/open/open_api/api/ws.api` | 路由定义 `GET /api/open/ws/v1/connect`（`appId` + `ticket`） |
| `app/open/open_api/internal/logic/ws/openwslogic.go` | 连接入口：鉴权 → 升级 → 注册 → 心跳 → 消息循环 |
| `app/open/open_api/internal/logic/ws/mqconsumerlogic.go` | 消费 `ws_bot_push_topic` 并推送 |
| `app/open/open_api/internal/logic/ws/websocket/enter.go` | 消息循环：控制帧（PING/PONG）+ 业务帧分派 |
| `app/open/open_api/internal/logic/ws/websocket/conn/client.go` | 连接封装（写互斥），复用 `ws_api/response` 与 `ws_api/types` |
| `app/open/open_api/internal/logic/ws/websocket/conn/manager.go` | `BotOnlineWsMap` 注册表 + `SendMsgToRobot` |
| `app/open/open_api/internal/logic/ws/websocket/auth/auth.go` | 机器人验签（`appId` + `ticket`） |
| `app/open/open_api/internal/logic/ws/websocket/heartbeat/manager.go` | 心跳（**不含 `coreonline` 调用**） |
| `app/open/open_api/internal/logic/ws/websocket/handler/chat_message/` | `enter.go` 按 `data.type` 分派 / `common.go` 收件人解析（**含会话归属校验**）/ `stream_send.go` 流式增量转发 |
| `common/wsEnum/wsTypeConst/wsTypeConst.go` | 新增 `ChatMessageStreamSend`、`ChatMessageStreamReceive` |
| `app/open/open_api/internal/handler/ws/openwshandler.go` | handler（把 `w/r` 透传给 logic 以支持升级） |
| `app/open/open_api/internal/svc/servicecontext.go` | 加 `RocketMQ` 客户端 |
| `app/open/open_api/open.go` | 启动事件消费者 |
| `app/open/open_api/etc/open.yaml` | 加 `RocketMQ` + `WebSocket` 配置段 |
| `app/open/openevent/publisher.go` | 新增 `Push()` —— 事件生产侧的**统一入口**，保证 MQ payload 形状一致 |
| `app/chat/chat_rpc/internal/logic/robot_event_push.go` | 由 `robot_webhook_push.go` 改名（机制早已不是 webhook）；`dispatch()` 由「调空实现 RPC」改为「投 MQ」；去掉失效的开关判断 |
| `app/open/open_api/internal/logic/robot/addrobottogrouplogic.go`、`removerobotfromgrouplogic.go` | 进群 / 出群事件由空 RPC 改为 `openevent.Push` |
| `app/friend/friend_api/internal/logic/deletefriendlogic.go`、`uservalidstatuslogic.go` | 取关 / 关注事件由空 RPC 改为 `openevent.Push` |

### 插件（beaver-sdk/packages/beaver-openclaw）

| 文件 | 说明 |
|---|---|
| `src/channel/long-connection.ts` | `BOT_WS_PATH` → `/api/open/ws/v1/connect`；新增 `sendStreamDelta()` —— 上行流式增量 |
| `src/accounts/index.ts` | 默认 `connectionMode = 'channel'` |

---

## 七、协议（服务端与插件的契约）

```
建连   GET {ws|wss}://<host>/api/open/ws/v1/connect?appId=..&ticket=..
凭据   ticket 复用 OAPI accessToken（POST /api/open/auth_public/v1/token）

业务帧（服务端 → 客户端）
  { "command": "CHAT_MESSAGE",
    "content": { "timestamp": 0, "messageId": "",
                 "data": { "type": "im.message.receive",
                           "conversationId": "conv_x",
                           "body": { ... } } } }
  data.type 承载事件类型（im.message.receive / im.message.receive.group_at）
  body 内容：sender_id / conversation_id / message_id / msg_type / content / mentions
  对齐 app/ws/ws_api/types/ws.go 的 WsMessage / WsContent / WsData

控制帧（双向）
  { "command": "PING" | "PONG" | "ACK", "messageId": "", "timestamp": 0 }

心跳   客户端每 25s 发 PING，服务端回 PONG
       服务端另有协议级 WebSocket Ping 帧（PingPeriod，默认 240s）
```

### 鉴权链（机器人独立分支）

`ticket 有效` → `ticket 归属该 appId` → `应用启用` → `应用已开启机器人能力（OpenRobot.Status=1）` → 得到 `robotId` 并注册连接。

与用户侧 `ws_api` 的 `VerifyWsToken` 完全分离：机器人没有用户登录态，那条链路依赖 Redis 的 `user_authentication_session:*`，不可复用。

---

## 八、关键设计决策

1. **放在 `open_api` 而非独立服务** —— 见第二节；网关零改动、零新增服务，隔离靠「独立注册表 + 不碰 coreonline + 独立 topic」实现。
2. **独立连接注册表 `BotOnlineWsMap`** —— 不复用 `UserOnlineWsMap`；索引用 `robotId`，命中是 O(1)，不是前缀扫描。
3. **机器人不上报用户在线态** —— 不调 `coreonline.MarkOnline`，后台在线统计与 `push_notify` 不受影响。
4. **独立 MQ Topic** —— 用户侧 `ws_api` 完全看不到机器人事件，零干扰。
5. **事件只投 MQ，不做 Webhook 回落** —— 平台的事件订阅表与投递逻辑已整体移除（`DispatchPlatformEvent` 现为空实现），Webhook 不再是可用通道。
6. **离线不补推** —— 消费者发现本实例没有该机器人的连接就直接丢弃。长连接模式下「没连上」即「对方离线」；离线期间的消息由机器人重连后自行拉取（**该拉取接口尚未实现，见第十节**）。
7. **复用 `CHAT_MESSAGE` 命令 + `ws_api/types` 帧结构** —— 不新增 WS 命令枚举，插件侧也更容易对齐 OpenClaw 标准事件。
8. **单机器人单连接（重复建连踢旧）** —— 注册表 key 是 `robotId`，不带设备维度。同一 `appId` 再连一次会踢掉旧连接。理由与「将来要多连接该怎么做」见第十二节。
9. **流式增量不落库，终稿才落库** —— 增量帧走旁路（`chat_message_stream_send` → `open_api` 转发 → `chat_message_stream_receive`），不进 `chat_rpc`、不占 seq、可丢帧；终稿仍走 OAPI `send_message` 落库。已实现，见第十一节。

---

## 九、风险与取舍

| 风险 | 说明 | 应对 |
|---|---|---|
| 多实例下连接归属 | 机器人连在实例 A，事件却可能在实例 B 产生 | MQ **广播消费**：每个实例都消费，谁持有连接谁推 |
| 企业代理拦 WebSocket | 部分企业网络会拦 ws | 目前无降级通道（Webhook 已移除） |
| 断线丢事件 | 长连接断开期间的事件不补推 | 待做：机器人重连后主动拉取（第十节） |
| 与开放平台 HTTP API 同进程 | 机器人连接数大时会拖累 `open_api` | 观察连接量级；需要时搬迁为独立服务 |
| go-zero 超时中间件 | 默认 `Timeout=3000ms` 是否会掐断 WS | **不会**：go-zero 的 `TimeoutHandler` 对 `Upgrade: websocket` 显式放行 |

---

## 十、待办 / 已知限制

1. **插件侧尚未按事件类型分派** —— 插件的 `handleBusinessFrame` 把所有帧都当「消息」解析（`validateMessage` 强制要求 `conversationId` / `senderId` / `content` 三者非空）。因此除 `im.message.receive*` 之外的 4 类生命周期事件（关注 / 取关 / 进群 / 出群）虽已投递到 MQ 并送达插件，但会被校验拦下、只打一条警告日志。**服务端已具备完整的事件下行能力，缺口在插件侧。**
2. **机器人「读消息历史」接口未实现** —— 这是「离线不补推」的前提。当前机器人只能在连接存活期间收到事件，断线期间的会话无法补齐。
3. **`DispatchPlatformEvent` 已成死代码** —— 所有调用方均已改为走长连接，现在只剩 proto 定义 + logic + server 注册。彻底删除需动 `open_rpc.proto` 并重新生成 pb，属跨模块清理，未做。
4. **`beaver/app/agent/agent_models` 包缺失** —— `main.go`、`database/migrations.go` 引用它，导致这两个包编译失败。**预先存在的问题**，与本次改动无关。
5. **端到端未实测** —— 需起 RocketMQ + `open_api`，建应用并开启机器人能力，插件用 `appId`/`appSecret` 连接后验证闭环。
6. **流式推送：服务端与插件已完成，客户端未做** —— 服务端转发链路（`handler/chat_message/`）与插件上行（`sendStreamDelta()`）已实现并编译通过，见第十一节。**缺口在三端客户端的渲染**：不实现客户端渲染的话，增量帧到达后会因客户端不认识该 `type` 而被忽略，退回"只有终稿"的效果（即做法 A），不会出错。
7. **机器人上行帧仍无错误回执** —— `HandleRobotMessages` 的 `default` 分支与 `chat_message.Handle` 的 `default` 分支现在都会以 **error 级别**带上具体 `command` / `type` 打日志（不再静默），但仍不会回帧给龙虾。彻底解决需要给 `WsControlFrame` 加 `code` 字段，而该结构体是与 `ws_api` 共用的，改动会影响用户侧，故未做。

---

## 十一、流式推送（龙虾 → IM → 客户端）

### 1. 现状：这条链路**不存在**，而且不是"没接上"，是"IM 没有这个概念"

三处证据：

| 位置 | 事实 |
|---|---|
| 插件 `messaging/outbound/outbound.ts` | `sendBeaverText` 一次性调 `POST /api/open/robot/v1/send_message`，把整段文本发完。OpenClaw 内部的流式分片在插件里就被 `StreamBuffer` 吃掉了，对 IM 不可见 |
| 服务端所有下行帧 | 全部是「落库后的完整快照」。`sendmsglogic.notifyMessageUpdateGrouped` 推的是 `body.tableUpdates = [messagesUpdate, conversationsUpdate, userConversationsUpdate]`，客户端按表覆盖，**没有"局部更新一条消息的某个字段"的语义** |
| `ChatMessage.Status` | 枚举里写了 `3=已编辑`，但**全仓没有任何代码路径会写它**。`chat_api` 也没有编辑消息内容的接口（只有 `recall` / `deleteMessages` / `markMessageMedia` / `forward` / `searchMessages`） |

唯一的先例是 `markMessageMedia`：它推 `chat_message_media_receive`，body 带 `{conversationId, messageIds}`，客户端在本地给已有消息打标。**这是"对已有消息做局部更新"的现成模式，流式可以照抄。**

### 2. 三种做法

| 做法 | 服务端改动 | 客户端改动 | 效果 | 结论 |
|---|---|---|---|---|
| **A. 不流式** | 无 | 无 | 龙虾攒完整段，一次性发出 | 插件现在就是这么做的，**今天就能跑** |
| **B. 主路流式**（先建消息，再反复 PATCH 同一条） | 新增 `update_message` 接口 + 新推送类型；高频写 DB、占 seq | 三端都要改 | 真流式 | **不推荐**：把流式写进 `ChatMessage` 会把 DB 和 MQ 打爆，且"半截消息"会被离线拉取、搜索、预览全部读到 |
| **C. 旁路流式**（增量走旁路，终稿走主路） | 新增 1 个上行命令 + 1 个转发逻辑 + 1 个推送类型；**不碰 chat_rpc / DB / seq** | 三端都要改 | 真流式 | **推荐** |

### 3. 推荐方案：旁路流式

```
增量（高频、可丢、不落库）
  龙虾 --WS 帧 STREAM_DELTA--> open_api --ws_push_topic--> ws_api --> 客户端

终稿（一次、必达、落库）
  龙虾 --OAPI send_message--> chat_rpc 落库 --ws_push_topic--> ws_api --> 客户端
```

客户端用龙虾生成的 `streamId` 把增量与终稿串起来：先按 `streamId` 起一个气泡，逐帧追加；收到终稿（带 `messageId`）后把气泡"定稿"为真实消息。

**为什么旁路是对的：**

1. **IM 核心零改动** —— 不碰 `chat_rpc`、不碰 `ChatMessage` 表、不碰 seq 分配。流式是"传输态"，不是"存储态"。
2. **增量丢帧无害** —— 下一帧会补齐，终稿兜底。不需要 ACK、不需要重传。
3. **历史记录天然正确** —— 只有终稿进 DB，离线拉取 / 搜索 / 会话预览读到的都是完整消息，不会读到半截。
4. **不污染 `coreonline`、不污染 `ws_push_topic` 的语义** —— 只是多了一种 `type`。

**落地情况（4 处）：**

| # | 改动 | 状态 |
|---|---|---|
| 1 | `websocket/enter.go` —— `HandleRobotMessages` 增加 `CHAT_MESSAGE` 业务帧分派 | ✅ 已完成 |
| 2 | `websocket/handler/chat_message/` —— 按 `data.type` 分派 + 收件人解析（**含会话归属校验**）+ 流式转发 | ✅ 已完成 |
| 3 | `wsTypeConst` —— 新增 `ChatMessageStreamSend` / `ChatMessageStreamReceive` | ✅ 已完成 |
| 4 | 三端客户端（Flutter / 桌面 / Web）—— 该类型的处理与气泡"打字机"渲染 | ⬜ **待做，需你定渲染方案** |
| 5 | 插件 —— `BeaverChannelRuntime.sendStreamDelta()` 供龙虾上行增量 | ✅ 已完成 |

第 2 步的鉴权是必须的：不校验会话归属，任何已接入的机器人都能给任意用户刷屏。实现放在 `common.go` 的 `resolveRecipients()`，与 ws_api 的 `getTypingPeerIDs` 同构（私聊取对端 / 群聊与圈子读 `chat_user_conversations`），但**多了一步「robotID 是否在成员里」的判断**。用表而不是 `GroupRpc`，是为了不引入新 RPC 依赖、且一张表同时覆盖群聊与圈子。

**没有新增 WS 命令枚举**：复用 `CHAT_MESSAGE` + `data.type` 分派，与用户侧的 `typing_send` / `typing_receive` 完全同构。

**帧格式（已实现，客户端按此对接）：**

上行（龙虾 → 服务端），插件方法 `sendStreamDelta()`：

```json
{ "command": "CHAT_MESSAGE",
  "content": { "timestamp": 0, "messageId": "",
    "data": { "type": "chat_message_stream_send",
              "conversationId": "private_a_b",
              "body": { "streamId": "st_xxx", "delta": "你", "seq": 3, "done": false } } } }
```

下行（服务端 → 客户端），由 `open_api` 转发：

```json
{ "command": "CHAT_MESSAGE",
  "content": { "timestamp": 0, "messageId": "",
    "data": { "type": "chat_message_stream_receive",
              "conversationId": "private_a_b",
              "body": { "streamId": "st_xxx",     // 龙虾生成，用于关联同一轮
                        "senderId": "<robotId>",
                        "delta": "你",             // 增量片段（不是全量）
                        "seq": 3,                  // 帧序号，客户端据此丢弃乱序帧
                        "done": false } } } }
```

客户端对接要点：

1. 按 `streamId` 建一个临时气泡，逐帧把 `delta` **追加**到内容尾部。
2. `seq` 跳号 → 丢弃该帧（说明有丢帧，等下一帧或终稿兜底）。
3. 收到 `chat_conversation_message_receive`（终稿，带真实 `messageId`）时，把临时气泡"定稿"为真实消息。**终稿是唯一权威**，增量只是草稿。
4. `done: true` 的增量帧不落库、不可回放——离线用户不会看到打字过程，只会看到终稿，这是符合预期的。

`delta` 用增量而非全量，是为了避免每帧都传整个长文本（长回复下带宽是平方级）。

---

## 十二、多连接：现状与建议

### 现状

`BotOnlineWsMap` 的 key 是 `robotId`（不含设备维度），`Register(robotID, client)` 在发现同一 `robotId` 已有连接时会**踢掉旧连接**并返回被踢数量。所以：**同一 `appId` 同一时刻只有一条连接。**

### 我的判断：**保持单连接是对的，不要急着改成多连接**

理由：

1. **事件语义是"一条消息 = 一个任务"，不需要多播。** 用户侧 `ws_api` 用 `userID_deviceType` 做多连接，是因为"我手机和电脑都要看到同一条消息"。机器人没有这个需求——它只是要把任务处理掉。
2. **多连接立刻引入"这条事件推给谁"的问题。** 当前 MQ 是广播消费，每个实例都收到；一旦同一 `robotId` 在多个实例上各有一条连接，就会**重复推送、重复回复**。要解决就必须在消费侧做一致性哈希（按 `conversationId`），复杂度陡增。
3. **踢旧连接其实是有用的脑裂保护。** 它保证同一时刻只有一个龙虾实例在收消息。运维上"新实例起来、旧实例被踢"也是符合直觉的滚动发布行为。
4. **龙虾要并发，应该在它自己进程内做 worker pool**，而不是靠多条 WS 连接。连接数是网络资源，worker 数是计算资源，两件事不该绑在一起。

### 如果将来确实要多连接（多 region / 多租户 / 单机连接数瓶颈）

正确做法不是"放开踢旧"，而是这四步：

1. 建连 URL 增加 `connectionId`（由龙虾实例生成），注册表 key 改为 `robotId:connectionId`。
2. 取消"踢旧"，改为**每机器人连接数上限** + 心跳超时清理。
3. MQ 消费侧按 `hash(conversationId) % N` 选连接——保证**同一会话的事件永远落到同一条连接**，会话内顺序不破。
4. 幂等仍然必须做：广播消费 + 重连窗口仍可能重复投递，龙虾侧按 `event_id` 去重（插件的 `core/idempotency.ts` 已有）。

---

## 十三、明确不做的事

- 不改动 `app/ws/ws_api`（用户侧连接、心跳、MQ 消费全部保持原样）
- 不改动 `core/coreonline`
- 不复用 `ws_push_topic`（避免用户侧消费者空跑）
- 不新增 WS 命令枚举（复用 `CHAT_MESSAGE`）
- 不给 agent 特权通道（走与外部机器人相同的开放平台机制）
- 流式不做客户端渲染（服务端与插件已具备能力，三端渲染独立排期，见第十一节）
- 本轮不放开多连接（理由见第十二节）
