# 压测套件

针对 beaver-server 的压测工具集，支持单接口压测和混搭链路压测。

## 目录结构

```
scripts/stress/
├── .env.example / .env       # 本地配置（各服务地址、Redis 密码），.env 已 gitignore
├── lib/                        # 公共层
│   ├── config.go              # 读 .env，各服务地址集中管理
│   └── http.go                # 公共 HTTP 工具
├── seed/                       # 造数据层（所有压测共用产出）
│   ├── users.go               # 批量注册测试账号 → accounts.json
│   ├── tokens.go              # 批量登录拿 token，回填 accounts.json
│   ├── friends.go             # 造好友关系（TODO）
│   └── conversations.go       # 造会话（TODO）
├── accounts.json              # 共用账号数据（gitignore，造数据后生成）
├── single/                     # 单接口压测（按服务分目录）
│   ├── auth/login.js          # 压 auth_api 登录
│   ├── chat/send_msg.js        # 压 chat_api 发消息
│   ├── chat/msg_list.js        # 压 chat_api 拉消息列表（TODO）
│   └── friend/...             # 压 friend_api 各接口（TODO）
├── mixed/                      # 混搭压测（模拟真实用户链路，跨服务）
│   └── im_flow.js             # 登录→拉好友→发消息→拉列表（TODO，待 seed 完成）
└── report/report.go           # JSON 报告转 HTML
```

## 设计思路

- **单接口压测**：看单个接口性能上限，按服务分目录，端口隔离
- **混搭压测**：模拟真实用户行为链路，多个接口按比例混合，看整体瓶颈
- **数据共用**：seed 造一次数据，single 和 mixed 都复用 `accounts.json`
- **配置隔离**：`.env` 存所有服务地址和密码，gitignore 永不进 git
- **可扩展**：加新接口就在对应服务目录加 `.js`，加新服务就建目录

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

### 3. 造数据
```powershell
cd seed
go run users.go -count 500    # 造账号
go run tokens.go              # 拿 token
```

## 压测流程

### 单接口压测
```powershell
cd single/auth
$env:AUTH_API_BASE="http://你的IP:20100"
k6 run --out json=report.json login.js
```

### 生成 HTML 报告
```powershell
cd ../../report
go run report.go ../single/auth/report.json ../single/auth/report.html
# 浏览器打开 report.html
```

### 混搭压测（待 seed 全部完成后）
```powershell
cd mixed
$env:AUTH_API_BASE="http://你的IP:20100"
$env:CHAT_API_BASE="http://你的IP:20300"
k6 run --out json=report.json im_flow.js
```

## 怎么看结果

- **http_req_duration**：延迟，p(99) < 1000ms 算及格
- **http_req_failed**：错误率，< 1% 算及格
- **iterations/s**：QPS
- 哪档并发延迟飙升/错误率上升 = 系统承载上限

## 扩展指南

- **加新单接口压测**：在 `single/<服务>/` 下新建 `.js`，照着 `auth/login.js` 改
- **加新服务的压测**：在 `single/` 下建目录，加 `.js`
- **加新造数据**：在 `seed/` 下新建 `.go`，用 `lib.Load()` 读配置
- **加新混搭场景**：在 `mixed/` 下新建 `.js`

## 状态

- [x] 骨架 + 配置层 + 造用户/造 token
- [x] auth_api 登录单接口压测
- [ ] seed：造好友关系、造会话
- [ ] chat_api 发消息、拉消息列表压测
- [ ] friend_api 好友列表、申请好友压测
- [ ] mixed：IM 完整链路混搭压测
