# 🦫 Beaver IM - 企业级即时通讯平台

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-2.1.1-blue.svg)](VERSION)
[![Go](https://img.shields.io/badge/go-1.24+-blue.svg)](https://golang.org/)
[![Go-Zero](https://img.shields.io/badge/Go--Zero-v1.7.4-green.svg)](https://github.com/zeromicro/go-zero)
[![QQ群](https://img.shields.io/badge/QQ群-1013328597%2B-blue.svg)](https://qm.qq.com/q/82rbf7QBzO)

> **海狸 IM 后端服务** - 基于 Go-Zero 微服务架构构建，为移动端（Flutter）、桌面端（Electron）与后台管理系统提供 REST / WebSocket / gRPC 能力。

**当前版本：[2.1.1](VERSION)**（以仓库根目录 [`VERSION`](VERSION) 文件为准）

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

本项目基于 [MIT](LICENSE) 协议开源，详见 [LICENSE](LICENSE)。

**使用要点（摘要）：**

- 闭源自用商用、二次开源均可免费，但须保留根目录 `LICENSE`，上线前端须有「关于」署名（基于海狸 IM + 仓库地址）
- 闭源交付第三方、去掉署名、对外 SaaS 收费等，请采购商业授权（书面合同）
- 无论是否付费，**不得删除或篡改 `LICENSE`**

完整免责与署名要求：[LEGAL.md](LEGAL.md)  
商业授权产品线与报价：[版权与商业授权](https://wsrh8888.github.io/beaver-docs/community/license.html)  
联系：[751135385@qq.com](mailto:751135385@qq.com)

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
