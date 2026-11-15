# go-zero 服务端（beaver-server）

> 合并自原 `beaver-server-dev` skill 的 api-conventions / architecture / gotchas 三份，并内嵌 `.cursor/rules/beaver-server.mdc`。

## 三条红线（违反必被否决）

1. 本仓 rpc 只在本仓服务间用，**绝不出仓**；仓内服务互调用 rpc，不要改成 S2S。
2. `beaver-server-package/` 下 6 个团队仓（ai/backend/call/doc/open/social）各自 rpc 不出业务线。
3. 跨仓一律走 **S2S**：打对方正常业务 API + Bearer 机器票。不造 `/internal/*` 门面，不加转发层。跨仓客户端叫 remote / S2S，不叫 rpc。

历史上曾把 `friend_api` 调用误改成 S2S、曾建 datasync 全量门面，均被要求回滚。

## 代码生成（必须遵守）

- 只改 `.api` / `.proto` / `logic`，用 goctl 生成 handler、routes、types。
- 禁止手改 `routes.go` / handler 等生成文件；禁止删除已有 handler / logic 文件。
- 禁止跨服务直接访问其他模块数据库（微服务隔离）。
- 禁止 `go build` / 启动服务 / 打包——除非用户明确要求。
- 改行为请改 `template/` 后重新生成，不要手改生成文件。
- Res 不含 `code`/`msg`，中间件包裹为 `{ code, msg, result }`。

## .api 文件与路由

- 每个 `@handler` 独立 `type()` 块（入参 + 出参 + 专用子类型），接口之间不共享请求/响应类型；字段注释写清含义与类型。
- 有 `@server group`、由主 `.api` import 子 `.api` 的服务用 **5 段**：`/{prefix}/{service}/{group}/v1/{action}`，如 `/api/open/auth_public/v1/token`、`/admin/backend/user/v1/list`。
- 无分组单模块用 **4 段**：`/{prefix}/{module}/v1/{action}`，如 `/api/friend/v1/update_notice`。
- `v1` 紧挨动作；`service` 是 service 名、`group` 是 `@server group`；管理后台固定 `/admin/backend/{group}/v1/{action}`。
- 必须含 `v1`；ID 用 `form` / `query`，禁止 `:id` 路径参数；**禁止 6 段及以上**。

### 消灭 6 段

6 段一律违规，"资源 + 操作"把资源并入动作名（snake_case）：
```
/api/auth/auth_public/v1/qrcode/status   → /api/auth/auth_public/v1/qrcode_status
/api/auth/auth/v1/qrcode/scan            → /api/auth/auth/v1/qrcode_scan
```
历史 3 条 qrcode 6 段路由、2 条 org 路由 v1 位置偏后（`/api/org/v1/department/tree` → `department_tree`）是现状偏差，不要照抄；改 URL 属破坏性变更须同步调用方并与用户确认。

## handler / logic 分工与命名

- handler：入参校验（非空、格式、范围），失败直接返回。
- logic：业务逻辑（DB、RPC、Redis）。
- 命名：`.api` 驼峰；`.rpc` / 表名 snake_case。

## 架构与调用方向

| 场景 | 方式 |
|---|---|
| 本仓服务 A → B | etcd gRPC，走 `client/<子域>` 窄接口包 |
| 本仓 → 收费仓（circle/doc/moment/open/ai） | `remote/<x>` 出站客户端（S2S） |
| 收费仓 → 本仓核心能力 | S2S 打本仓业务 API（account/message/relation/asset/platform/notification） |
| 收费仓 → 本仓同步能力 | datasync 仅 `getSync*` |

### 服务速查（合并后 9 个）

| 服务 | 合并前 | api/rpc 端口 | etcd key | 库 |
|---|---|---|---|---|
| account | user+auth | 20100/30100 | accountrpc.rpc | beaver_account |
| relation | friend+group | 20200/30200 | relationrpc.rpc | beaver_relation |
| message | chat | 20300/30300 | messagerpc.rpc | beaver_message |
| asset | file+emoji | 20600/30600 | assetrpc.rpc | beaver_asset |
| notification | — | 21000 | notificationrpc.rpc | — |
| platform | — | 21600 | platformrpc.rpc | — |
| datasync | — | 21700 | datasyncrpc.rpc | — |
| ws / gateway | — | 20400 / — | — | — |

库名合并映射（表名不变）：friend/group→relation、file/emoji→asset、user/auth→account、chat→message。

### 合并范式（goctl multiple）

单一 proto（package `xxx_rpc`，go_package `./types`），生成：
```bash
goctl rpc protoc xxx.proto --go_out=. --go-grpc_out=. -m --zrpc_out=.
```
产物：`types/` 单 pb 包、`client/<子域>/`、`internal/{logic,server}/{子域}/`；单 etcd 注册；已存在的 etc/svc/config/main 不被覆盖。

- 调用方 pb 别名：account→`accttypes.`、asset→`assettypes.`、relation→`types.`（冲突 alias `relationtypes`）、message→`chat_rpc.`。
- client 构造：`NewUser` / `NewAuth` / `NewFriend` / `NewGroup` / `NewFile` / `NewEmoji`。
- gateway `ServiceAlias`（URL 前缀→合并服务名，7 条）：friend/group→relation、file/emoji→asset、user/auth→account、chat→message。deploy yaml 的 `Name` 不带 `_api` 后缀（代码注册 `c.Name+"_api"`）。根 Dockerfile = 9 服务 16 二进制，改名后同步。

## 鉴权（双通道）

同一套 `/api` 路由同时服务用户 JWT 与 S2S 机器票，差异只在 `middleware/authgate/`。机器票单受众 `core:api`，换票端点 `POST /api/auth/token`（account_api，client_credentials）。鉴权失败返回真实 401（非 200+code=1，否则 S2S 客户端无法清票重换）。

新服务接入：config 加 `Auth` + `S2S`、svc 构造 `authgate.NewGate(...)` 并 `Validate()`、main 里 `server.Use(ctx.AuthGate)`。**不要另造转发层或第二套门面。**

特例：agent_api 用静态 `X-Internal-Token`（非 OAuth），与调用方 yaml `InternalToken` 同值。

## go.mod replace（本地联调）

- `beaver-ai/server`：`replace beaver => ../../../beaver-server`
- `beaver-backend` / `beaver-open`：`replace beaver-ai => ../beaver-ai/server`

模块名注意：`beaver-social` 的 module 是 `github.com/wsrh8888/beaver-server-social`；跨仓复制代码时按目标仓 go.mod 实际 module 改 import 前缀。

## 踩坑清单

### CRLF / 批量编辑

- 仓库文件多为 **CRLF**：python/sed 的 `\n` 匹配会静默失败。用 `splitlines(keepends=True)` 行过滤、正则配 `\r?\n`、跳过块终止条件用 `line.rstrip() == "}"`。
- bash 传 `python -c` 时 `\n` 被损坏成字面 `/n`：换行用 `chr(10)` 或写文件。sed 替换侧 `&` 需转义 `\&`。
- MSYS 路径 `/f/...` Windows Python 打不开，用 `F:/...`；`python /dev/stdin` 不可用。
- 收费仓在 `F:/code/mine/github/beaver-server-package/` 下；本仓 import 是短名 `beaver/...`，grep 别用全 module 名（假阴性）。

### goctl

- api 生成：`goctl api go -api xxx.api -dir . --home ../../../template`（本仓自带模板，response 信封由模板 + 中间件统一包裹，不手改 handler）。
- `.api` 不接受 C 风格块注释 `/** */`（用 `//`）；同文件 service 名唯一；新增 service 名可能误建额外 main 文件（删本次误生成的，不删已有）。
- `duplicate type expression` 的根因是接口间共享类型——按规范每个 `@handler` 独立 `type()` 块。
- rpc 生成后检查 pb 落盘目录（由 `go_package` 决定）。

### pb 类型映射

`.api` types 与 proto types 是两套 Go 类型，用 `common/convert.To[T]`（JSON 往返）转换，字段名 / json tag 需逐字一致；转换失败当业务错误抛出，不要静默返回零值。

常见名不一致：`.api` `UserLookup/UserList/UserUpdate/UserDelete/UserUpdateStatus/UserSearch/ResetUserPassword` → proto `UserInfo/ListUsers/UpdateUsers/DeleteUsers/UpdateUsersStatus/SearchUser/ResetPassword`。新增 DTO 先对照 `beaver-backend/types/remote/server/*.go` 的 json tag。

### 版权头

所有 Go 文件保留完整 `beaver-server-header` 版权头，生成的新文件也要补。
