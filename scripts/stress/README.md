# 压测套件

针对 beaver-server 的压测工具集，支持单接口压测和混搭链路压测。
观测走标准 k6 → InfluxDB → Grafana 时序栈，看延迟/QPS 随并发变化曲线。

## 目录结构

```
scripts/stress/
├── .env.example / .env        # 本地配置（各服务地址、Redis 密码），.env 已 gitignore
├── docker-compose.yml          # InfluxDB + Grafana 可观测性栈
├── observability/grafana/      # Grafana 数据源自动注入
├── utils/                     # 基础设施层（每个关注点一个独立 Go package）
│   ├── config/               # package config：读 .env，AccountsPath()/FriendsPath() 统一数据路径
│   │   └── config.go
│   ├── httputil/             # package httputil：HTTP 客户端（全局复用连接池，不叫 http 避免遮蔽 stdlib）
│   │   └── http.go
│   ├── types/                # package types：纯数据类型 Account / Friendship
│   │   └── account.go
│   └── k6/                    # JS flow 层（k6 脚本共用，必须独立，JS 无包系统）
│       ├── config.js         # 从 .env 读服务地址（Go/k6 共用一份配置）
│       ├── shared.js         # 账号池/think/请求头/结果校验
│       ├── auth.js           # k6 flow：login
│       ├── friend.js        # k6 flow：addFriend / friendList / friendInfo
│       └── chat.js          # k6 flow：sendMsg / recentList
├── setup/                     # 建状态层（自带 Go flow，唯一消费者就是 setup，不独立成层）
│   ├── users/main.go         # 造账号 → accounts.json（含 registerOne flow）
│   ├── tokens/main.go        # 造 token，回填 accounts.json（含 loginOne flow）
│   └── friends/main.go       # 配对互加好友 → friends.json（含 makeFriends/addFriend/valid/friendInfo flow）
├── accounts.json              # 共用账号数据（gitignore）
├── friends.json               # 好友关系 + 会话ID（gitignore，setup/friends 产出）
├── single/                     # 单接口压测（按服务分目录，调 k6 flow）
│   ├── auth/login.js         # 压 auth_api 登录
│   ├── chat/send_msg.js      # 压 chat_api 发消息（conversationId 取自 friends.json）
│   └── friend/list.js        # 压 friend_api 好友列表
└── scenarios/                  # 场景压测（按场景分目录，每场景独立，调 k6 flow 串链路）
    └── daily_chat/daily_chat.js  # 日常聊天：拉好友→发消息→拉会话列表
```

## 设计思路：utils 基础设施 + setup 自带 flow + k6 共用 flow

- **utils/ = 基础设施**：配置/HTTP/数据类型，纯 plumbing，不含业务。Go flow 不放这里。
- **setup/ = 建状态，自带 Go flow**：每个 setup 文件内联自己的真实业务 flow（registerOne/loginOne/makeFriends 等），Go flow 唯一消费者就是 setup，不独立成层；worker pool + 限速编排。
- **utils/k6/ = JS flow 共用层**：k6 脚本（single + scenarios）跨脚本复用同一份 flow，JS 无包系统必须独立成目录。
- **single/ + scenarios/ = 调 k6 flow 跑负载**：单接口看上限，场景看链路瓶颈；都调 utils/k6 里的 flow。
- **数据共用**：setup 造一次，single/scenarios 都复用 accounts.json + friends.json；路径用 `utils.AccountsPath()/FriendsPath()` 统一。
- **配置隔离**：.env 存所有服务地址和密码，gitignore 永不进 git。
- **标准可观测**：k6 → InfluxDB → Grafana，看时序曲线而非单次快照。

### 鉴权约定（基于各服务 main.go 中间件注册）

- **auth_api**：登录需要 UA 做设备识别（设备类型校验 + 互踢）
- **friend_api**：仅 UA 解析中间件，不校验 token，只需 `Beaver-User-Id` header
- **chat_api**：无中间件，只需 `Beaver-User-Id` header
- → 业务接口请求头不带 Authorization，`utils/k6/shared.js` 的 `userHeaders()` 只设 `Beaver-User-Id`

### 会话产生方式

会话随好友关系自动创建（`friend_info` 返回 conversationId），**无独立建会话接口**。
所以 setup 只需造好友关系，会话自动产生，conversationId 写入 friends.json 供 chat 压测取用。

## 前置准备

### 1. 安装 k6（本机）
https://github.com/grafana/k6/releases 下 Windows 版，解压把 `k6.exe` 放到 PATH。
验证：`k6 version`

### 2. 填写本地配置
```powershell
cd scripts/stress
copy .env.example .env
# 编辑器打开 .env 填写：各服务 IP:端口、Redis 密码
```

### 3. 起可观测性栈（一次性）
```powershell
cd scripts/stress
docker compose up -d           # 起 influxdb(8086) + grafana(3000)
```

### 4. 造数据
```powershell
cd setup
go run ./users -count 500 -c 10 -rps 5    # 造账号（并发=10 限速=5rps）
go run ./tokens -c 10                       # 拿 token（并发=10）
go run ./friends -c 5                        # 造好友关系 + 会话（并发=5，涉及双向请求不宜高）
```
- `-c` 并发 worker 数
- `-rps` 每秒提交上限（仅 users，保护服务端频率限制）

## 压测流程

服务地址统一从 .env 读（k6 脚本通过 utils/k6/config.js 自动加载），无需手动设环境变量。
每个脚本按自己所在目录算相对路径，单接口在 single/<服务>/ 下、场景在 scenarios/<场景名>/ 下执行。

### 单接口压测
```powershell
# auth 登录
cd single/auth
k6 run --out influxdb=http://localhost:8086/k6 login.js

# friend 好友列表
cd single/friend
k6 run --out influxdb=http://localhost:8086/k6 list.js

# chat 发消息
cd single/chat
k6 run --out influxdb=http://localhost:8086/k6 send_msg.js
```

### 场景压测
```powershell
cd scenarios/daily_chat
k6 run --out influxdb=http://localhost:8086/k6 daily_chat.js
```

## 看结果

### 控制台摘要（每次必看，零配置）
k6 跑完自动打印 `http_req_duration` 的 p(95)/p(99)、`http_req_failed` rate、`iterations/s`（QPS）、VU 数。
`thresholds` 不达标 k6 退出码非 0，可直接接 CI 门禁。

### 时序仪表盘（InfluxDB + Grafana，看趋势）
浏览器开 http://localhost:3000（admin/admin）→ Dashboards → Import → 填官方 k6 仪表盘 ID **2587**。

### 本地归档（离线分析）
```powershell
k6 run --out json=report.json login.js
```
原始数据点逐行 JSON，归档或离线分析用。

## 怎么判读

- **http_req_duration** p(99) < 1000ms 算及格
- **http_req_failed** < 1% 算及格
- 哪档并发延迟飙升/错误率上升 = 系统承载上限

## 扩展指南

- **加新接口**：Go 侧 flow 内联到对应 setup 文件（如需建状态）+ k6 侧 flow 加到 `utils/k6/<服务>.js`（压测共用）
  - 例：加群聊建状态 → `setup/groups/main.go` 内联 `createGroup` flow；加群聊压测 → `utils/k6/group.js` 的 `createGroup`
- **加新单接口压测**：在 `single/<服务>/` 下新建 `.js`，调 `utils/k6/` 里的 flow，照着 `friend/list.js` 改
- **加新建状态脚本**：在 `setup/<名字>/main.go` 里内联自己的 flow + worker pool 编排，照着 `friends/main.go` 改
- **加新场景压测**：在 `scenarios/<场景名>/` 下新建目录 + `.js` + `README.md`，串多个 flow，照着 `daily_chat/` 改

## 状态

- [x] 骨架 + utils 基础设施层 + 可观测性栈
- [x] setup 自带 Go flow：users(register) + tokens(login) + friends(makeFriends)
- [x] utils/k6 JS flow：auth + friend + chat（跨脚本共用）
- [x] single：auth/login + friend/list + chat/send_msg
- [x] scenarios：daily_chat 真链路（拉好友→发消息→拉会话列表）
- [ ] single：chat 会话列表、消息历史等读接口
- [ ] single：group/moment 等其他服务
