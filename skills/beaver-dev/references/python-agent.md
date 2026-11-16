# Python agent（beaver-ai）

> 仓根：`beaver-server-package/beaver-ai/`。  
> 跨仓红线与 go-zero 一致：见 `go-zero-server.md`（rpc 不出仓；跨仓只走 S2S / remote）。

## 环境与运行

在仓根执行（不要先 `cd` 进 `app/ai/agent_api` 再套相对路径）：

```bash
cd beaver-server-package/beaver-ai
uv sync
uv run python main.py -db                  # 按 ORM create_all 建表
uv run python app/ai/agent_api/main.py     # 唯一 HTTP 进程（接入 + 配置态 + 跑图）
```

- 本域 **无独立 agent_rpc**：原 rpc 能力已并入 `agent_api`，禁止再起第二进程或对本域打 gRPC。
- `-db` 只按当前 ORM `create_all`，不要在 `main.py` 里手写 `ALTER` 补列。本机旧表可删库 `beaver_agent` 再 `-db`；空库直接 `-db`。
- `common.paths.boot()`：`app` → `agent_api/app`，`beaver` → `agent_api/beaver`（运行时内核与 HTTP 同进程）。

## 分层（对标 go-zero api，禁止再造 domain）

```text
app/ai/
  agent_api/                 # 唯一 API 进程
    app/
      handler/               # 入参校验、组 Res
      logic/                 # 业务（DB / remote）
      middleware/            # HTTP 闸门（鉴权、请求日志）；单 API 不放仓根 middleware
      svc/                   # ServiceContext 装配
      types/                 # 本服务入出参
    beaver/                  # 跑图 / tools / runtime 内核（非 HTTP 层）
    etc/
  agent_models/              # 本域表模型（一文件一表，禁止塞进 agent_api/app）
```

- handler → logic → models/DB 或 `remote/<peer>`；**禁止** `app/domain`、`internal/domain`。
- 仓根 `core/`（coregorm 等）与 `beaver/core`（log/metrics）职责不同，不要强行合并。

## Remote / S2S（强制，对齐 `beaver-server/remote`）

跨仓出站 **只叫 remote / S2S，不叫 rpc**。结构与用法对齐 Go 仓 `beaver-server/remote` + `types/remote`，Python 落在 beaver-ai 仓根：

```text
beaver-ai/
  commonconfig/              # 出站 peer 域名与凭据注册（AiUrl/DocUrl/OpenUrl…）
  remote/
    <peer>/                  # 一个对端产品一个包：doc / open / beaver / …
      __init__.py            # 或 enter.py：Init(client_id, client_secret)
      remote.py              # Post/Get/Delete + 换票缓存（plumbing，不抽跨 peer 公共包）
      token.py               # GetTokenApi：打对端 /machine/v1/token
      <biz>.py               # 业务方法，每个对端能力一个函数
  types/
    remote/
      __init__.py
      token.py               # TokenRequest / TokenResponse 等共用换票类型
      <peer>/
        *.py                 # 该 peer 的 Req/Res dataclass（本仓自持，不 import 对端模块）
```

### 硬性规则

1. **只打对方正常业务 / machine HTTP**，六段路径；禁止 `/internal/*` 门面，禁止 etcd gRPC 出仓。
2. **Req/Res 只放 `types/remote/<peer>`**，本仓自维护；禁止依赖对端 Go/Python 业务包里的类型。
3. 每个 peer 包自带 `remote.py` plumbing（换票 + 带票发包），**不**抽「全仓统一 RemoteClient」把各线路缠在一起（与 Go remote 各业务线各存一份一致）。
4. `Init` 只装配本 peer 凭据（及必要时写入 `commonconfig` 域名）；业务侧直接 `remote.doc.get_doc_api(...)`，调用行不散落换票逻辑。
5. 出站鉴权：优先 `Authorization: Bearer`（client_credentials 换票）；静态 `X-Internal-Token` 仅对仍声明该方式的对端使用。
6. remote **不重试、不聚合、不缓存业务结果**；只缓存换票（含过期 skew）。
7. logic 里调 remote；handler 不直接调 remote。

### 方法命名（强制 `Api` 后缀）

出站业务函数名 **必须以 `Api` 结尾**，一眼区别于本仓 logic / 本地 DB 方法：

| 对 | 错 |
|---|---|
| `get_doc_api` | `get_doc` |
| `get_doc_content_api` | `fetch_doc_content` |
| `get_agent_api` | `get_agent` |
| `get_token_api` | `get_token` / `exchange`（对外导出换票也要 `Api`） |

约定：

- Python：`snake_case` + `_api` 后缀，例如 `get_doc_api(ctx, headers, body) -> Res`。
- 若将来 Go remote 同步改造：`GetDocApi` / `GetTokenApi`（PascalCase + `Api`）。
- `remote.py` 内的 `Post` / `Get` / `Delete` / `auth_header` 属于通道工具，**不加** `Api`。
- 函数只做「拼 `commonconfig.XxxUrl + 路径` + 调 Post/Get/Delete」；路径写全，不在调用方拼。

示例（语义对齐 `beaver-server/remote/doc/content.go`，命名带 Api）：

```python
# remote/doc/content.py
async def get_doc_api(ctx, headers, body: GetDocReq) -> GetDocRes:
    return await post(
        ctx,
        commonconfig.DOC_URL + "/beaver-doc/api/doc/machine/v1/get_doc",
        headers,
        body,
    )

async def get_doc_content_api(ctx, headers, body: GetDocContentReq) -> GetDocContentRes:
    return await post(
        ctx,
        commonconfig.DOC_URL + "/beaver-doc/api/doc/machine/v1/get_doc_content",
        headers,
        body,
    )
```

### 何时新建 peer 包

- 本仓要调外部产品（doc / open / beaver-server gateway 等）→ 建 `remote/<peer>` + `types/remote/<peer>`。
- 仅本进程内读 `agent_models` / 跑图 → **不要**包一层 remote。

## 工具通道

- 能力清单与本机执行在 `app/ai/agent_api/beaver/`（mcp / `tools/client_session.py`）。
- 事件名 `client_tool`；回传走 agent_api 已有 HTTP（如 `submit_client_tool_result`）。
- 不要再引入 host、sandbox。

## 库表与业务约定

- `agent_workspaces.path` 只给提示词和界面展示；真正写磁盘在拥有该工作空间的那台机器上。
- 用量写 `agent_usage_ledger`，消息表没有 `LlmMeta`。
- 消息必须有 `branch_id`，话题没有当前分支时建主分支。
- 在线状态是别的微服务的事，agent 不查设备是否在线。

## Middleware

- 放在 `agent_api/app/middleware/`（auth、request_log）。
- 公网面：采信网关 `Beaver-User-Id`。
- 机器面：Bearer 机器票；换票入口自身豁免；`X-Internal-Token` 仅作显式兜底。
- 验签/换票纯逻辑可放仓根 `utils/s2s/`，middleware 只做 HTTP 接线。

## 已知问题（未授权不要改）

桌面收到 `client_tool` 后未校验是否拥有该工作空间，两台电脑在线会都执行、先提交的生效——没让修之前不要改。

## 生成与版本

- 本域已无对外 gRPC；不要再为 agent 本域生成/提交 pb 桩。
- 若其它遗留工具仍用 grpc：仓根 `uv` 环境里的 `grpc_tools`（protobuf 4.x、grpcio 1.62），不要用机器全局更新的 protoc。
