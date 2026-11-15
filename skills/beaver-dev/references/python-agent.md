# Python agent（beaver-ai）

> 来源：`.cursor/rules/beaver-agent.mdc`。注意：该 `.mdc` 里的 `beaver-agent/agent_server` 是旧路径，实际仓根已迁到 `beaver-server-package/beaver-ai/`，以下命令用新路径。

## 环境与运行

环境在仓根 `beaver-server-package/beaver-ai`（这一层执行，不要先 `cd` 进 `app/ai/agent_api` 再套相对路径）：

```bash
cd beaver-server-package/beaver-ai
uv sync
uv run python main.py -db                  # 核心服务
uv run python app/ai/agent_rpc/main.py     # agent rpc
uv run python app/ai/agent_api/main.py     # agent api（静态 X-Internal-Token，端口 22300）
```

- `-db` 只按当前 ORM `create_all`，不要在 `main.py` 里手写 `ALTER` 补列。本机旧表删掉 `beaver_agent` 再跑 `-db`；换电脑空库直接 `-db`。
- `common.paths.boot()` 里 `app` 固定是 `agent_api/app`，`beaver` 在 `agent_rpc`。

## 工具通道

- 能力清单在 `app/ai/agent_rpc/beaver/mcp`，本机执行在 `beaver/tools/client_session.py`（完整路径 `app/ai/agent_rpc/beaver/tools/client_session.py`）。
- 事件名 `client_tool`，回传 `POST /api/agent/v1/message/submitClientToolResult`。
- 不要再引入 host、sandbox。

## 库表与业务约定

- `agent_workspaces.path` 只给提示词和界面展示；真正写磁盘在拥有该工作空间的那台机器上。
- 用量写 `agent_usage_ledger`，消息表没有 `LlmMeta`。
- 消息必须有 `branch_id`，话题没有当前分支时建主分支。
- 在线状态是别的微服务的事，agent 不查设备是否在线。

## 已知问题（未授权不要改）

桌面收到 `client_tool` 后未校验是否拥有该工作空间，两台电脑在线会都执行、先提交的生效——没让修之前不要改。

## 生成与版本

- Python 桩用仓根环境里的 `grpc_tools` 生成（protobuf 4.x、grpcio 1.62）。生成后的 gRPC 导入保持 `from agentpb import agent_rpc_pb2 as agent__rpc__pb2`，不要用机器上更新的 protoc。
- 长上下文参考：`.workbuddy-ai/memory/2026-10-11.md` 后半「同日后续」一节。
