# 🦫 Beaver IM - 企业级即时通讯平台

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-2.0.2-blue.svg)](VERSION)
[![Go](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
[![Go-Zero](https://img.shields.io/badge/Go--Zero-v1.7.4-green.svg)](https://github.com/zeromicro/go-zero)
[![QQ群](https://img.shields.io/badge/QQ群-1013328597%2B-blue.svg)](https://qm.qq.com/q/82rbf7QBzO)

> **海狸 IM 后端服务** - 基于 Go-Zero 微服务架构构建，为移动端（Flutter）、桌面端（Electron）与后台管理系统提供 REST / WebSocket / gRPC 能力。

**当前版本：[2.0.2](VERSION)**（以仓库根目录 [`VERSION`](VERSION) 文件为准）

[English](README_EN.md) | [中文](README.md)

---

## ✨ 核心能力

- 🔐 **认证授权** - 手机号 / 邮箱 / 扫码 / OAuth 登录，JWT 鉴权，登录设备管理
- 💬 **即时通讯** - 私聊与群聊，消息存储与检索，已读回执，WebSocket 实时推送
- 👥 **社交关系** - 好友申请与资料、群组创建与管理、朋友圈
- 📁 **文件与表情** - 文件上传下载、表情包与收藏
- 📞 **音视频通话** - 基于 LiveKit 的通话信令与房间管理
- 🔔 **消息通知** - 系统通知与互动提醒
- 🔄 **跨端同步** - 数据同步服务，支持多端消息与媒体状态一致
- 🧩 **开放平台** - 开放 API 与开发者门户，支持机器人 / Webhook 集成
- 📦 **平台服务** - 字典、意见反馈、版本更新、客户端日志查询等公共能力
- 🛠️ **后台管理** - 用户管理、消息审计、连接监控、版本发布等管理端接口

## 🛠️ 技术栈

| 类别 | 技术 | 说明 |
|------|------|------|
| 语言 / 框架 | Go 1.24 + Go-Zero 1.7 | 微服务 API / RPC |
| 通信 | gRPC、WebSocket、RocketMQ | 服务间调用与实时消息 |
| 存储 | MySQL 8.0、Redis 6.0 | 多库拆分 + 缓存 |
| 基础设施 | ETCD 3.5、Docker | 服务发现与部署 |
| 其他 | JWT、LiveKit、七牛云 | 鉴权、通话、对象存储 |

## 🏗️ 架构概览

客户端（Flutter / Electron / 管理端）统一经 API 网关访问各业务服务；服务间通过 gRPC 协作，WebSocket 负责长连接推送。

```
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Flutter 移动端 │  │ Electron 桌面端│  │  后台管理 Web │
└──────┬───────┘  └──────┬───────┘  └──────┬───────┘
       │                 │                 │
       └─────────────────┼─────────────────┘
                         │
              ┌──────────▼──────────┐
              │  gateway_api :20800  │
              │  gateway_admin :40800│
              └──────────┬──────────┘
                         │
    ┌────────────────────┼────────────────────┐
    │                    │                    │
 user/auth/friend   chat/ws/group/file    platform/open/call
 moment/emoji/...   notification/datasync   backend :40000
    │                    │                    │
    └────────────────────┼────────────────────┘
                         │
              ┌──────────▼──────────┐
              │ MySQL · Redis · ETCD │
              └─────────────────────┘
```

### 服务端口

| 服务 | API 端口 | RPC 端口 | 说明 |
|------|---------:|---------:|------|
| user | 20000 | 30000 | 用户 |
| auth | 20100 | 30100 | 认证 |
| friend | 20200 | 30200 | 好友 |
| chat | 20300 | 30300 | 聊天 |
| ws | 20400 | - | WebSocket |
| group | 20500 | 30500 | 群组 |
| file | 20600 | 30600 | 文件 |
| emoji | 20700 | 30700 | 表情 |
| gateway_api | 20800 | - | 客户端 API 网关 |
| moment | 20900 | 30800 | 朋友圈 |
| notification | 21000 | 31000 | 通知 |
| platform | 21600 | 31100 | 平台（字典 / 反馈 / 更新 / 日志等） |
| datasync | 21700 | - | 数据同步 |
| call | 21800 | 31800 | 通话 |
| open_api | 21900 | 30900 | 开放平台 API |
| open_portal | 22000 | - | 开放平台门户 |
| backend | 40000 | - | 后台管理 |
| gateway_admin | 40800 | - | 管理端网关 |

> 完整端口表见 [`server.md`](server.md)；实际部署以各服务 `etc/*.yaml` 为准。

## 📁 项目结构

```
beaver-server/
├── app/                    # 微服务
│   ├── auth/               # 认证
│   ├── user/               # 用户
│   ├── friend/             # 好友
│   ├── chat/               # 聊天
│   ├── ws/                 # WebSocket
│   ├── group/              # 群组
│   ├── file/               # 文件
│   ├── emoji/              # 表情
│   ├── moment/             # 朋友圈
│   ├── notification/       # 通知
│   ├── platform/           # 平台公共能力
│   ├── datasync/           # 跨端同步
│   ├── call/               # 音视频通话
│   ├── open/               # 开放平台
│   ├── backend/            # 后台管理
│   └── gateway/            # API 网关
├── common/                 # 中间件、响应、校验等公共组件
├── core/                   # 数据库 / Redis / ETCD 等核心配置
├── database/               # 初始化数据与迁移脚本
├── deploy/                 # 部署配置
├── main.go                 # 数据库创建与 AutoMigrate 入口
└── server.md               # 服务端口说明
```

## 🚀 快速开始

详细的环境准备、配置说明与部署步骤见文档站：

- 📖 [后端开发文档](https://wsrh8888.github.io/beaver-docs/backend/)
- 📖 [部署说明](https://wsrh8888.github.io/beaver-docs/backend/deploy/build-scripts.html)

本地初始化数据库（需先配置 MySQL 连接）：

```bash
go run main.go -db
```

## 📚 文档与资源

- 📖 **详细文档**: [https://wsrh8888.github.io/beaver-docs/](https://wsrh8888.github.io/beaver-docs/)
- 🎥 **视频教程**: [B站频道](https://space.bilibili.com/269553626/lists)
- 📱 **体验包下载**: [海狸 IM Android 体验包](https://github.com/wsrh8888/beaver-docs/releases/download/lastest/latest.apk)
- 💬 **QQ 群**: [1013328597](https://qm.qq.com/q/82rbf7QBzO)

## 🔗 相关项目

| 项目 | 仓库地址 | 说明 |
|------|----------|------|
| **beaver-server** | [GitHub](https://github.com/wsrh8888/beaver-server) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-server) | 后端微服务（本仓库） |
| **beaver-flutter** | [GitHub](https://github.com/wsrh8888/beaver-flutter) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-flutter) | 移动端（Flutter，推荐） |
| **beaver-desktop** | [GitHub](https://github.com/wsrh8888/beaver-desktop) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-desktop) | 桌面端（Electron） |
| **beaver-manager** | [GitHub](https://github.com/wsrh8888/beaver-manager) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-manager) | 后台管理系统 |
| **beaver-open** | [GitHub](https://github.com/wsrh8888/beaver-open) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-open) | 开放平台 |
| **beaver-oauth** | [GitHub](https://github.com/wsrh8888/beaver-oauth) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-oauth) | OAuth 授权登录 |
## 🤝 贡献指南

欢迎通过 Issue / Pull Request 参与贡献，详见 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 📄 开源协议与免责声明

本项目基于 [MIT](LICENSE) 协议开源 - 详情请参阅 [LICENSE](LICENSE) 文件。

### ⚖️ 使用说明

**项目定位**：本项目主要用于**技术学习和交流**，希望为开发者提供一个学习和研究的平台。

**使用建议**：
- 📚 **学习交流** - 欢迎用于个人学习、技术研究、学术交流
- 🤝 **开源贡献** - 欢迎提交代码改进、Bug修复、功能建议
- 🔒 **合规使用** - 请确保使用方式符合当地法律法规
- 💡 **创新应用** - 鼓励基于本项目进行创新性应用开发

**温馨提示**：
- 本项目采用 MIT 开源协议，您可以自由使用、修改和分发
- 建议在使用前仔细阅读相关法律法规，确保合规使用
- 如有疑问或需要帮助，欢迎通过 QQ 群或 GitHub Issues 交流

### 📋 项目来源标注要求

**重要**：如果您基于本项目进行二次开发或发布，**必须**在项目中保留以下信息：

#### 🖥️ **前端项目（移动端/桌面端/Web应用）**
- **关于页面**：必须在"关于我们"、"关于应用"或类似页面中包含项目来源标注
- **必需文本**："基于 [Beaver IM](https://github.com/wsrh8888/beaver-server) 开源IM项目开发"
- **链接**：必须提供可点击的原始项目链接

#### 🔧 **后端项目（服务器/API服务）**
- **README.md**：必须在项目介绍或描述中包含来源标注
- **必需文本**："基于 [Beaver IM](https://github.com/wsrh8888/beaver-server) 开源IM项目开发"
- **链接**：必须提供可点击的原始项目链接

#### 📄 **通用要求**
- **LICENSE 文件**：保留原项目 MIT 协议信息

> 💡 **友好提醒**：本项目允许个人及商业使用；基于本项目二次开发或发布时，**必须保留项目来源标注**，详见上方要求。

> 📖 **详细法律条款**：请参阅 [LEGAL.md](LEGAL.md) 文件了解完整的法律免责声明和使用要求。

## ⭐ Star 历史

[![Star History Chart](https://api.star-history.com/svg?repos=wsrh8888/beaver-server&type=Date)](https://star-history.com/#wsrh8888/beaver-server&Date)

## ☕ 请作者喝杯茶

如果这个项目对你有帮助，欢迎请作者喝杯茶 ☕

<div align="center">
  <img src="./static/sponsor/wechat.jpg" width="200" alt="微信赞助码"/>
  <img src="./static/sponsor/zhifubao.jpg" width="200" alt="支付宝赞助码"/>
</div>

---

<div align="center">
  <strong>Made with ❤️ by Beaver IM Team</strong><br>
  <em>企业级即时通讯平台</em>
</div>
