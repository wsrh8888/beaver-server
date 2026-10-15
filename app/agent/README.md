# app/agent

beaver-server 侧的 agent 模块仅是 **RPC 客户端桩**（proto + 生成代码），不跑 go-zero 服务端。

## 目录结构

- `agent_rpc/agent_rpc.proto` — gRPC 服务定义（与 beaver-agent 仓库的 `agent_server/grpc/agent/agent_rpc.proto` 双向对齐，权威定义在 beaver-agent 侧）。
- `agent_rpc/agent/agent.go` — goctl 生成的 zrpc 客户端 `agent.Agent`（`agent.NewAgent(zrpc.MustNewClient(...))`）。
- `agent_rpc/types/agent_rpc/` — protoc 生成的 Go 消息类型。

## 用途

供 beaver-server 各模块（如 `chat_rpc` 的 agent_task / agent_done 编排层）通过 gRPC 调用 **beaver-agent 仓库** 中的 Python `agent_server` 服务（`SendMessage`、`SubmitClientToolResult`、`InvalidateModel`、`ReportCapabilities`、`ListConversations`、`ListTopics`、`ListMessages`、`UpsertWorkspace`）。

## 不跑服务端

本目录**不启动** go-zero 服务端。权威的 agent_rpc 服务端在 beaver-agent 仓库（Python 实现），通过 etcd 注册为 `agentrpc.rpc`。本目录只消费该服务。

## 重新生成

proto 变更后，在 `agent_rpc/` 目录下执行（proto 文件末尾注释即该命令）：

```bash
goctl rpc protoc agent_rpc.proto --go_out=./types --go-grpc_out=./types --zrpc_out=.
```

注意：`--zrpc_out=.` 会顺带生成一份 go-zero 服务端脚手架（`agentrpc.go` + `internal/` + `etc/`），本仓库**不需要** Go 侧服务端，生成后请删除这些脚手架文件，仅保留 `agent/`（客户端）与 `types/`（消息类型）。
