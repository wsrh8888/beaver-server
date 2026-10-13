# 场景：日常聊天（daily_chat）

模拟一个真实用户在日常聊天时的完整操作链路，带并发负载跑，看跨服务链路瓶颈。

## 链路

```
拉好友列表（friend_api）→ 发消息（chat_api）→ 拉最近会话列表（chat_api）
```

| 步骤 | 接口 | 服务 | 读写 | 说明 |
|------|------|------|------|------|
| 1 | `GET /api/friend/v1/friend_list` | friend_api | 读 | 看自己有哪些好友 |
| 2 | `POST /api/chat/v1/sendMsg` | chat_api | 写 | 给好友发一条文本消息 |
| 3 | `GET /api/chat/v1/getRecentChatList` | chat_api | 读 | 看最近会话列表（发消息后会刷新） |

跨服务：friend_api + chat_api。登录在 setup 阶段已拿过 token，链路里不重复登录。

## 为什么测这个场景

单接口压测看的是**单个接口的性能上限**，但真实用户不会只调一个接口——他是一连串操作。串起来带负载跑，能暴露单接口看不出的瓶颈：

- **连接池竞争**：friend 和 chat 两个服务同时被压，RPC/DB 连接池是否够
- **DB 行锁/热点**：发消息会更新会话最后一条消息，拉列表又读会话，并发下会话行的读写竞争
- **缓存一致性**：发消息后会话列表该刷新，缓存失效/重建在高并发下的表现
- **token 校验开销**：每个请求都带 Beaver-User-Id，鉴权解析是否成热点

这些只有链路 + 负载才能暴露，单接口测不出来。

## flow 复用

每个接口调 `utils/k6/` 里的 flow，与单接口压测共用同一份 flow，不在脚本里重写 HTTP 细节：

- `utils/k6/friend.js` 的 `friendList`
- `utils/k6/chat.js` 的 `sendMsg` / `recentList`

`conversationId` 从 `setup/friends` 产出的 `friends.json` 取（会话随好友关系自动创建，无独立建会话接口）。

## 前置

setup 全部跑完：

```powershell
cd setup
go run ./users -count 500 -c 10 -rps 5    # 造账号
go run ./tokens -c 10                        # 拿 token
go run ./friends -c 5                         # 造好友关系 + 会话
```

## 运行

```powershell
cd scenarios/daily_chat
$env:CHAT_API_BASE="http://你的IP:20300"
$env:FRIEND_API_BASE="http://你的IP:20200"
k6 run --out influxdb=http://localhost:8086/k6 daily_chat.js
```

## 负载模型

`ramping-vus` 阶梯加压，看哪档并发开始出瓶颈：

| 阶段 | 时长 | 目标 VU |
|------|------|---------|
| 升压 | 2m | 50 |
| 保持 | 3m | 50 |
| 升压 | 2m | 100 |
| 保持 | 3m | 100 |
| 降压 | 1m | 0 |

## 看什么

- **http_req_duration p(99) < 1500ms**（链路比单接口宽，因为是跨服务串行）
- **http_req_failed < 2%**
- 在 Grafana（仪表盘 ID 2587）看延迟/QPS 随 VU 变化曲线，定位拐点
- 哪档并发延迟飙升或错误率上升 = 链路承载上限
