# go-zero 服务端（beaver-server）

## 公共（全仓）

### 三条红线

1. 本仓 rpc 只在本仓服务间用，**绝不出仓**；仓内互调走 rpc，不要改成 S2S。
2. `beaver-server-package/` 下各团队仓（ai/backend/call/doc/open/social）各自 rpc 不出业务线。
3. 跨仓一律 **S2S**：打对方业务 API + Bearer 机器票。不造 `/internal/*` 门面。跨仓客户端叫 remote / S2S，**不叫 rpc**。

### 通用硬性

- 禁止跨服务直接访问其他模块数据库。
- 禁止 `go build` / 启动服务 / 打包——除非用户明确要求。
- 禁止删除已有 handler / logic 文件；废弃接口先停声明再谈。
- 只改契约与业务实现：`.api` / `.proto` / `logic`（及确需的 `svc` / `config` / `template`）；生成物用 goctl，**禁止手改** `routes.go` / handler / pb 生成文件等。
- 改公共生成行为：改 `template/` 后重新生成。
- 新 Go 文件保留完整 `beaver-server-header`；生成文件也要补。
- 命名：`.api` 驼峰；`.proto` / 表名 snake_case。
- handler：入参校验；logic：业务（DB / rpc / remote）。

### 调用方向

| 场景 | 方式 |
|---|---|
| 本仓服务 A → B | etcd gRPC（rpc） |
| 本仓 → 独立产品仓 | `remote/<x>`（S2S） |
| 独立仓 → 本仓 | S2S 打本仓业务 API |

### 领域目录（强制）

每个业务域目录下**只允许**这三类（命名 `xx_api` / `xx_models` / `xx_rpc`）：

```text
app/<domain>/
  <domain>_api/       # HTTP
  <domain>_models/    # 本域库表模型
  <domain>_rpc/       # 仓内 gRPC（可按域裁剪，但不得另起杂目录替代 models）
```

示例：`account_api` + `account_models` + `account_rpc`；`open_api` + `open_models` + `open_rpc`。  
禁止在域下再塞 `app/`、`common_models/`、随意 `models/` 等第三套表目录。

---

## Models（数据库表）规范

目录按**业务**分子包，**一文件一表**：

```text
xx_models/
  user/
    user_model.go              # 仅 user 表
    user_setting_model.go      # 仅 user_setting 表
  auth/
    auth_credential_model.go
    auth_device_model.go
  oauth/
    open_oauth_model.go        # 仅 open_oauth 表
    open_oauth_token_model.go
```

硬性：

- 表模型必须落在本域 `xx_models`，禁止写进 `_api` / `_rpc` 的 internal。
- **按业务分子目录**（`user` / `auth` / `oauth` / `bot` / `digital`…），目录名用业务语义，**禁止**含糊的 `app` 之类大杂烩包。
- **一个 `.go` 文件只对应一张表**（一个 `TableName()` / 一个主 model struct）。禁止 1 个文件里定义多张表。
- 该表 JSON 列的 `Scan`/`Value` 辅助类型可放同文件，且**只服务这一张表**；不能顺手配别的表。
- 文件名：`xxx_model.go`（或现有 `*_model.go` 风格），表名 snake_case。

反例（不合法）：

- `open_models/app/open_oauth_model.go`：业务目录叫 `app` 不合法；OAuth 表应在 `open_models/oauth/`，且该文件不得再夹其它表。
- 一个 `xxx_models.go` 里堆 `User`、`Order`、`Token` 多张表。

---

## API 服务规范

目录形态（强制）：

```text
xxx_api/
  xxx_api.api          # 主入口：只 import，不写业务 type/handler
  api/                 # 必须有；按 group 拆子 .api
    auth_public.api
    user.api
    ...
  internal/
    handler/           # goctl 生成
    logic/             # 业务手写
    types/             # goctl 生成
    ...
  xxx.go               # main
```

主文件只做引入，例如：

```api
import "api/auth_public.api"
import "api/auth.api"
import "api/user.api"

// goctl api go -api xxx_api.api -dir . --home ../../../template
```

### 路由（六段）

```text
/{team}/{surface}/{service}/{group}/v1/{action}
```

| 段 | 含义 |
|---|---|
| team | `beaver` 或 `beaver-open` 等独立仓 |
| surface | `api` / `admin` / 机器面（如 `machine`） |
| service | 服务名，如 `auth`、`open` |
| group | 与 `@server group`、子 `.api` 一致 |
| v1 | 固定，紧挨 action |
| action | snake_case；资源+操作并进动作名 |

示例：`/beaver/api/auth/auth_public/v1/phone_login`  
禁止：`:id` 路径参数；禁止 `/v1/qrcode/status` 再拆段（写成 `qrcode_status`）。

网关只转发、不鉴权；鉴权在下游 `@server` 声明（`jwt` / `middleware`）；公开组不挂身份中间件。鉴权失败真实 **401**。

### `.api` 契约

- **每一个** `@handler` = 一个接口 = **一组**独立 `type ( ... )`。
- 该 `type` 块内必须有本接口的 `XxxReq` + `XxxRes`（Res 可空：`XxxRes {}`），以及**仅本接口**用的子类型。
- **不允许接口间共享任何 type**（含 Req、Res、嵌套子结构）。每个接口契约独立；共享会导致改 A 连带改坏 B。字段看起来一样也要各自拷一份、各自命名（如 `PhoneLoginUser` / `EmailLoginUser`），禁止共用 `User` / `CommonReq` 之类。
- 有 N 个 `@handler` 就有且仅有 N 个业务 `type` 块。
- Res **不含** `code` / `msg`；中间件包成 `{ code, msg, result }`。
- 只用 `//` 注释；同文件 `service` 名唯一；字段注释写清含义。

```api
type (
	PhoneLoginReq {
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	PhoneLoginRes {
		Token  string `json:"token"`
		UserID string `json:"userId"`
	}
)

@server (
	group: auth_public
)
service account {
	@handler phoneLogin
	post /beaver/api/auth/auth_public/v1/phone_login (PhoneLoginReq) returns (PhoneLoginRes)
}
```

### API 代码生成

```bash
goctl api go -api xxx_api.api -dir . --home ../../../template
```

只改 `api/*.api`（及主 import）和 `internal/logic`；生成 handler / routes / types，不手改。

---

## RPC 服务规范

目录形态（强制，**multiple / `-m`**）：

```text
xxx_rpc/
  xxx.proto            # 单 proto，多个 service
  client/<子域>/       # goctl -m 生成窄客户端
  internal/
    logic/<子域>/
    server/<子域>/
    ...
  pb/                  # protobuf 生成物（go-zero 官方习惯名）
  xxx.go
```

### proto

- 一个进程一个 proto 包（如 `package account_rpc`）。
- **pb 落点（强制）**：`option go_package = "./pb";` → 生成到本服务根下 `pb/*.pb.go`（不要进 `internal/types`，也不要用旧的 `types/`）。
- **必须多个 `service`** 按子域拆（如 `service user`、`service auth`），编译加 **`-m`**；禁止单 service 把所有 RPC 糊成一团。
- 字段注释写清；命名 snake_case。

### RPC 代码生成

与 go-zero 官方一致，`go_package` / 生成目录都叫 `pb`：

```bash
goctl rpc protoc xxx.proto --go_out=. --go-grpc_out=. -m --zrpc_out=. --home ../../../template
```

- `option go_package = "./pb";`：goctl/protoc 会按包名自动建 `pb/` 子目录，所以 out 参数用 `.`（不是 `./pb`），产物恰好落在 `xxx_rpc/pb/*.pb.go`。实测当前 goctl 版本即如此；若写 `--go_out=./pb` 反而会套娃成 `pb/pb/`。
- `-m` 必带：产出 `client/<子域>`、`internal/{logic,server}/<子域>`。
- 只改 `.proto` 与对应 `logic`；`pb/` 下生成物、server、client 不手改。
- rpc **仅仓内**互调；外仓禁止直连本仓 etcd rpc。

### 与 API 类型

两套目录、两套类型，勿混：

| 来源 | 目录 |
|---|---|
| `.api` → goctl api | `xxx_api/internal/types/` |
| `.proto` → protoc | `xxx_rpc/pb/` |

跨层用 `common/convert.To[T]`，json tag 必须一致，失败当业务错误，不要静默零值。

---

## 踩坑

- 仓库多为 CRLF；路径用 `F:/...`。
- 发现两个接口引用同一 type，或 goctl 报 `duplicate type expression`：立刻拆开，禁止共享；原因是改 A 会影响 B。
- models 出现 `app/` 杂包或一文件多 `TableName`：按 Models 节拆到业务子目录，一文件一表。
- 新增 service 名可能误建额外 main：只删本次误生成的。
- 本仓 import 短名 `beaver/...`；独立仓按各自 go.mod module。
- rpc 生成物目录统一 `pb/` + `go_package = "./pb"`；改落点时 import 路径、无别名文件的 `types.` 限定符、生成文件 `package` 名必须一起迁移，不要半改半留。
