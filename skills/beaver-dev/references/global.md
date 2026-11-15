# 全局约定（跨端）

> 来源：workspace `.cursor/rules/rule.mdc` 的通用部分。

## 项目构成

开源 IM（Beaver IM），对标大厂 IM 产品。五端：

| 端 | 目录 | 技术栈 |
|---|---|---|
| 服务端 | `beaver-server` | go-zero |
| Electron 客户端 | `beaver-desktop` | Electron + Vite + Vue3 |
| Flutter 移动端 | `beaver-flutter` | Flutter 3 / Dart |
| Python agent | `beaver-server-package/beaver-ai` | Python（uv 管理） |
| Vue 后台 | `beaver-manager` | Vue3 + TS + Vite |

## Agent 改代码流程

- 涉及修改代码前，先说明要改哪些文件、改什么，等用户明确回复「可以改 / 确认 / 继续」后再调用编辑工具。
- 未确认前，只做阅读、搜索、分析，不直接改文件。
- 用户明确说「直接改」时可跳过确认。

## 通用硬性规则

- **不要写 md 文件**（skill 本身、用户明确要求的除外）。
- 新建或改文件，顶部版权头与现有文件一致、完整保留整段（含中英文说明），禁止缩短、删改、替换。版本标记按端区分：服务端 `beaver-server-header`，PC 端 `beaver-desktop-header-v2`。
- 文件命名不以大写字母开头，多个单词用驼峰命名。

## 架构总览

服务端 `beaver-server` 是核心 IM，通过网关对各端提供 HTTP API；各端不直连服务端数据库。跨团队/跨仓能力统一走 S2S（client_credentials 换票 + Bearer 访问业务 API），不走内部 RPC。详细调用方向见 `go-zero-server.md`。
