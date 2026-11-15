# 架构与调用方向

## 谁调谁

| 场景 | 方式 |
|---|---|
| 本仓服务 A → B | etcd gRPC，走 `client/<子域>` 窄接口包 |
| 本仓 → 收费仓（circle/doc/moment/open/ai） | `remote/<x>` 出站客户端（S2S） |
| 收费仓 → 本仓核心能力 | S2S 打本仓业务 API（account/message/relation/asset/platform/notification） |
| 收费仓 → 本仓同步能力 | datasync 仅 `getSync*` |

**红线**：本仓 rpc 不出仓；每个团队 rpc 不出业务线；跨仓一律 S2S，不造 `/internal/*` 门面、不加转发层。历史上曾把 `friend_api` 调用误改成 S2S、曾建 datasync 全量门面，均被要求回滚。

跨仓客户端叫 **remote / S2S**，不叫 rpc（「不能调别人团队的 rpc，只能用人家提供的 API」）。

## 服务速查（合并后 9 个）

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

## 合并范式（goctl multiple）

单一 proto（package `xxx_rpc`，go_package `./types`），生成：
```bash
goctl rpc protoc xxx.proto --go_out=. --go-grpc_out=. -m --zrpc_out=.
```
产物：`types/` 单 pb 包、`client/<子域>/`、`internal/{logic,server}/{子域}/`；单 etcd 注册。已存在的 etc/svc/config/main 不被覆盖。

调用方 pb 别名：account→`accttypes.`、asset→`assettypes.`、relation→`types.`（冲突 alias `relationtypes`）、message→`chat_rpc.`。client 构造：`NewUser`/`NewAuth`/`NewFriend`/`NewGroup`/`NewFile`/`NewEmoji`。

gateway `ServiceAlias`（URL 前缀→合并服务名，7 条）：friend/group→relation、file/emoji→asset、user/auth→account、chat→message。deploy yaml 的 `Name` 不带 `_api` 后缀（代码注册 `c.Name+"_api"`）。根 Dockerfile = 9 服务 16 二进制，改名后同步。

## 鉴权（双通道）

同一套 `/api` 路由同时服务用户 JWT 与 S2S 机器票，差异只在 `middleware/authgate/`。机器票单受众 `core:api`，换票端点 `POST /api/auth/token`（account_api，client_credentials）。鉴权失败返回真实 401（非 200+code=1）。新服务接入：config 加 `Auth`+`S2S`、svc 构造 `authgate.NewGate(...)` 并 `Validate()`、main 里 `server.Use(ctx.AuthGate)`。**不要另造转发层或第二套门面。**

## go.mod replace（本地联调）

- `beaver-ai/server`：`replace beaver => ../../../beaver-server`
- `beaver-backend`/`beaver-open`：`replace beaver-ai => ../beaver-ai/server`

模块名注意：`beaver-social` 的 module 是 `github.com/wsrh8888/beaver-server-social`；跨仓复制代码时按目标仓 go.mod 实际 module 改 import 前缀。
