---
name: beaver-server-dev
description: Beaver IM 核心服务 beaver-server 的开发导航与架构红线。当改动 .api/.proto/logic、服务合并、跨仓 S2S 调用、鉴权、goctl 代码生成、端口/etcd 配置时使用。提供「rpc 不出仓」等硬约束、跨仓调用方向、踩坑清单。
agent_created: true
---

# Beaver Server 开发导航

## 三条红线（违反必被否决）

1. 本仓 rpc 只在本仓服务间用，**绝不出仓**；仓内服务互调就用 rpc，不要改成 S2S。
2. `beaver-server-package/` 下 6 个团队仓，各自 rpc 不出业务线。
3. 跨仓一律走 **S2S**：打对方正常业务 API + Bearer 机器票。不造 `/internal/*` 门面，不加转发层。

## 先读哪个

- **改 `.api` / `.proto` / logic** → `references/api-conventions.md`（API 与路由规范，内嵌自 workspace `.cursor/rules/beaver-server.mdc`）。
- **改调用链**（谁调谁、rpc 还是 S2S）→ `references/architecture.md`。
- **动手前**（CRLF、goctl、类型映射）→ `references/gotchas.md`。

硬性禁止与路由段数规则同时记在 `api-conventions.md`，写接口前先看。

## 硬性禁止（项目规则）

- 禁止手改 `routes.go` / handler 等生成文件：只改 `.api` / `.proto` / logic，其余用 goctl + 本仓 `template/` 生成。
- 禁止删除已有 handler / logic 文件；禁止跨服务直连对方数据库。
- 禁止擅自 `go build` / 启动服务 / 打包——除非用户明确要求。
- Res 不含 `code`/`msg`（中间件统一包裹为 `{code, msg, result}`）。

## 改动前

Go 文件保留完整 `beaver-server-header` 版权头。批量编辑先读 gotchas（仓库多为 CRLF，常规 `\n` 匹配会静默失败）。
