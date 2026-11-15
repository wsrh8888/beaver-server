# API 与路由规范

> 副本来源：workspace 规则 `F:/code/mine/github/.cursor/rules/beaver-server.mdc`（`description: beaver-server go-zero API 与路由规范`）。此处内嵌以让 skill 自包含；若两边不一致，以 `.mdc` 为准并同步本文件。

## 代码生成（必须遵守）

- 只改 `.api` / `.proto` / `logic`，用 goctl 生成 handler、routes、types。
- 禁止手改 `routes.go` / handler 等生成文件。
- 禁止删除已有 handler / logic 文件。
- 禁止跨服务直接访问其他模块数据库（微服务隔离）。
- 禁止 `go build` / 启动服务 / 打包——除非用户明确要求。
- Res 不含 `code`/`msg`，中间件已包裹为 `{ code, msg, result }`。

## .api 文件

- 每个 `@handler` 对应独立 `type()` 块（入参 + 出参 + 该接口专用子类型）。
- 接口之间不共享请求/响应类型。
- 字段注释写清含义与数据类型。

## API 路由

有 `@server group`、由主 `.api` import 子 `.api` 的服务用 **5 段**；没有分组的单模块用 **4 段**。`v1` 紧挨动作。

- **5 段**：`/{prefix}/{service}/{group}/v1/{action}`
  ```
  /api/open/auth_public/v1/token
  /admin/backend/user/v1/list
  ```
  `service` 是 service 名，`group` 是 `@server group`；管理后台固定 `/admin/backend/{group}/v1/{action}`。同一 group 里的资源写进动作名，不要再加路径段。
- **4 段**：`/{prefix}/{module}/v1/{action}`
  ```
  /api/friend/v1/update_notice
  /api/user/v1/get_info
  ```
- 通用：必须含 `v1`；ID 用 `form` / `query`，禁止 `:id` 路径参数；**禁止 6 段及以上**。

### 段数决策：以 5 段为准，消灭 6 段

6 段一律违规。"资源 + 操作"把资源并入动作名（snake_case/驼峰），不新增路径段：

```
# 错 → 对
/api/auth/auth_public/v1/qrcode/status   → /api/auth/auth_public/v1/qrcode_status
/api/auth/auth_public/v1/qrcode/generate → /api/auth/auth_public/v1/qrcode_generate
/api/auth/auth/v1/qrcode/scan            → /api/auth/auth/v1/qrcode_scan
```

理由：少一层 = 匹配更直接、gateway 前缀映射更稳、与全仓 44 条 5 段路由一致；REST 语义靠动作名前缀表达即可。

**现状偏差（不要照抄）**：历史代码有 3 条 6 段 qrcode 路由（即上列），另有 2 条 org 路由 v1 位置偏后（`/api/org/v1/department/tree`、`/api/org/v1/member/list`，应压平为 `department_tree` / `member_list`）。改 URL 属破坏性变更，须同步调用方并与用户确认。

## handler / logic 分工

- handler：入参校验（非空、格式、范围），失败直接返回。
- logic：业务逻辑（DB、RPC、Redis）。

## 命名

- `.api`：驼峰；`.rpc` / 表名：snake_case。
