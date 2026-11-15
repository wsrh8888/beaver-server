---
name: beaver-dev
description: Beaver IM 全栈开发导航与项目铁律。覆盖五个技术栈——go-zero 服务端(beaver-server)、Python agent(beaver-server-package/beaver-ai)、Electron 客户端(beaver-desktop)、Flutter 移动端(beaver-flutter)、Vue 后台(beaver-manager)。当改任一端代码、查架构红线（rpc 不出仓 / 跨仓 S2S）、API 路由规范、goctl 代码生成、鉴权、UI 组件或项目通用约定时使用。已内嵌 workspace .cursor/rules 全部规范，多工具通用。
agent_created: true
---

# Beaver 开发导航

Beaver IM 五端统一开发规范。改代码前先看**全局铁律**，再按技术栈读对应 reference。

## 全局铁律（跨端，先看 `references/global.md`）

1. **改代码前先说明改哪些文件、改什么，等用户明确「可以改 / 确认 / 继续」再动手**；用户说「直接改」可跳过。
2. 不写 md 文件（skill 本身、用户明确要求的除外）。
3. 新建 / 改文件，版权头与现有文件一致、完整保留（`beaver-server-header` / `beaver-desktop-header-v2` 等），禁止缩短删改。
4. 文件命名不以大写开头，多单词用驼峰。

## 按技术栈读哪个

| 你要改 | 目录 | reference |
|---|---|---|
| go-zero 服务端 | `beaver-server` | `references/go-zero-server.md` |
| Electron 客户端 | `beaver-desktop` | `references/electron-desktop.md` |
| Flutter 移动端 | `beaver-flutter` | `references/flutter-mobile.md` |
| Python agent | `beaver-server-package/beaver-ai` | `references/python-agent.md` |
| Vue 后台 | `beaver-manager` | `references/vue-admin.md` |

通用约定、架构总览、改代码流程 → `references/global.md`。

## 服务端架构红线（最高优先级）

- 本仓 rpc 只在本仓服务间用，**绝不出仓**；仓内互调就用 rpc。
- `beaver-server-package/` 下 6 个团队仓各自 rpc 不出业务线。
- 跨仓一律 S2S：打对方正常业务 API + Bearer 机器票；不造 `/internal/*` 门面、不加转发层。详见 `go-zero-server.md`。

## 多工具安装

主体为纯 Markdown，一份内容装到任意工具：

- **WorkBuddy / Claude**：整个 `beaver-dev/` 放到 `~/.workbuddy-ai/skills/`（或建软链），自动识别 `SKILL.md`。
- **Cursor**：把 `cursor/beaver-dev.mdc` 复制到项目 `.cursor/rules/`；它已内嵌核心铁律可独立工作，`references/` 放在项目 `beaver-server/skills/beaver-dev/` 下供按需读取。
- **其他编辑器**：把 `SKILL.md` 作为项目指令文件读入即可。
