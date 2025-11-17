# 🦫 Beaver IM - 企业级即时通讯平台

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/)
[![Go-Zero](https://img.shields.io/badge/Go--Zero-v1.6.0+-green.svg)](https://github.com/zeromicro/go-zero)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()
[![Coverage](https://img.shields.io/badge/coverage-85%25-brightgreen.svg)]()
[![Stars](https://img.shields.io/badge/stars-50+-yellow.svg)](https://github.com/wsrh8888/beaver-server/stargazers)

> 🚀 **企业级即时通讯平台** - 基于Go-Zero微服务架构构建，支持移动端（UniApp）、桌面端（Electron）和Web端，提供实时通信能力。

[English](README_EN.md) | [中文](README.md)

---

## 🌟 核心特性

### 🔐 **企业级安全**
- **多因子认证** - 邮箱验证、短信验证码、生物识别支持
- **端到端加密** - 消息传输和存储加密
- **基于角色的访问控制** - 细粒度权限管理和管理员控制
- **审计日志** - 全面的安全事件追踪

### 💬 **高级消息功能**
- **实时通信** - 基于WebSocket的即时消息
- **多格式支持** - 文本、图片、文件、语音消息、表情
- **消息状态** - 已读回执、正在输入提示、送达确认
- **消息搜索** - 全文本对话搜索
- **消息撤回** - 限时消息删除

### 👥 **社交功能**
- **联系人管理** - 二维码扫描、联系人导入导出
- **群组管理** - 创建、管理和审核群聊
- **好友请求** - 带自定义消息的审批流程
- **用户资料** - 丰富的个人资料信息和头像

### 🏗️ **微服务架构**
- **15+微服务** - 可扩展、可维护的服务分解
- **高可用性** - 多实例部署支持
- **服务发现** - 基于ETCD的服务注册
- **负载均衡** - 智能请求分发
- **熔断器** - 故障容忍和弹性

### 📱 **跨平台支持**
- **移动应用** - 通过UniApp框架支持iOS/Android
- **桌面应用** - 通过Electron支持Windows/macOS/Linux
- **API网关** - 第三方集成的RESTful API

## 🛠️ 技术栈

### 后端服务
| 技术 | 版本 | 用途 |
|------|------|------|
| **Go-Zero** | v1.6.0+ | 微服务框架 |
| **gRPC** | v1.58+ | 服务间通信 |
| **WebSocket** | - | 实时消息 |
| **MySQL** | 8.0+ | 主数据存储 |
| **Redis** | 6.0+ | 缓存和会话管理 |
| **ETCD** | 3.5+ | 服务发现和配置 |
| **Docker** | 20.0+ | 容器化 |

### 前端技术
| 平台 | 框架 | 特性 |
|------|------|------|
| **移动端** | UniApp + Vue 3 | 跨平台移动应用 |
| **桌面端** | Electron + Vue 3 | 原生桌面体验 |

## 📊 性能指标

- **消息延迟**: 平均 < 100ms
- **并发用户**: 支持 10,000+
- **消息吞吐量**: 100,000+ 消息/秒
- **可用性**: 99.9% 正常运行时间
- **响应时间**: API响应 < 200ms

## 🏗️ 架构概览

```
┌─────────────────┐    ┌─────────────────┐
│   移动端应用     │    │   桌面端应用     │
│   (UniApp)      │    │   (Electron)    │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────────────────┘
                    │
                    ┌─────────────┴─────────────┐
                    │      API网关              │
                    │      (端口: 20800)        │
                    └─────────────┬─────────────┘
                                  │
        ┌─────────────────────────┼─────────────────────────┐
        │                         │                         │
┌───────▼────────┐    ┌───────────▼──────────┐    ┌────────▼────────┐
│   认证服务      │    │    用户服务           │    │   好友服务      │
│   API:20100    │    │   API:20000          │    │  API:20200      │
│   RPC:30100    │    │   RPC:30000          │    │  RPC:30200      │
└────────────────┘    └──────────────────────┘    └─────────────────┘
        │                         │                         │
┌───────▼────────┐    ┌───────────▼──────────┐    ┌────────▼────────┐
│   聊天服务      │    │    群组服务          │    │   文件服务      │
│   API:20300    │    │   API:20500          │    │  API:20600      │
│   RPC:30300    │    │   RPC:30500          │    │  RPC:30600      │
└────────────────┘    └──────────────────────┘    └─────────────────┘
        │                         │                         │
┌───────▼────────┐    ┌───────────▼──────────┐    ┌────────▼────────┐
│   WS服务       │    │   表情服务           │    │ 反馈服务        │
│   API:20400    │    │   API:20700          │    │  API:21400      │
│   -            │    │   -                  │    │  -              │
└────────────────┘    └──────────────────────┘    └─────────────────┘
        │
        └─────────────────────────────────────────────────────────────┐
                                                                      │
                    ┌─────────────────────────────────────────────────┴─┐
                    │              数据层                               │
                    │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
                    │  │    MySQL    │  │    Redis    │  │    ETCD     │ │
                    │  │   (8.0+)    │  │   (6.0+)    │  │   (3.5+)    │ │
                    │  └─────────────┘  └─────────────┘  └─────────────┘ │
                    └───────────────────────────────────────────────────┘
```



## 📚 文档与资源

- 📖 **详细文档**: [https://wsrh8888.github.io/beaver-docs/](https://wsrh8888.github.io/beaver-docs/)
- 🎥 **视频教程**: [B站频道](https://space.bilibili.com/269553626/lists)
- 📱 **体验包下载**: [海狸IM Android体验包](https://github.com/wsrh8888/beaver-docs/releases/download/lastest/latest.apk)
- 💬 **QQ群**: [1013328597](https://qm.qq.com/q/82rbf7QBzO)

## 🔗 相关项目

| 项目 | 仓库地址 | 说明 |
|------|----------|------|
| **beaver-server** | [GitHub](https://github.com/wsrh8888/beaver-server) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-server) | 后端微服务 |
| **beaver-mobile** | [GitHub](https://github.com/wsrh8888/beaver-mobile) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-mobile) | 移动端应用 |
| **beaver-desktop** | [GitHub](https://github.com/wsrh8888/beaver-desktop) \| [Gitee](https://gitee.com/dawwdadfrf/beaver-desktop) | 桌面端应用 |



## 📱 功能展示

### 🔐 用户认证
<div align="center">
  <img src="./static/mobile/login.jpg" width="200" alt="登录界面"/>
  <img src="./static/mobile/register.jpg" width="200" alt="注册界面"/>
  <img src="./static/mobile/find-password.jpg" width="200" alt="找回密码"/>
</div>

### 💬 聊天功能
<div align="center">
  <img src="./static/mobile/message.jpg" width="200" alt="消息主界面"/>
  <img src="./static/mobile/private-chat.jpg" width="200" alt="私聊聊天"/>
  <img src="./static/mobile/group-chat.jpg" width="200" alt="群聊聊天"/>
  <img src="./static/mobile/send-text.jpg" width="200" alt="发送文字"/>
  <img src="./static/mobile/send-emoji.jpg" width="200" alt="发送表情"/>
  <img src="./static/mobile/chat-details.jpg" width="200" alt="聊天详情"/>
</div>

### 👥 社交功能
<div align="center">
  <img src="./static/mobile/friend.jpg" width="200" alt="好友列表"/>
  <img src="./static/mobile/new-friends.jpg" width="200" alt="新的朋友"/>
  <img src="./static/mobile/friend-info.jpg" width="200" alt="好友资料"/>
  <img src="./static/mobile/edit-remark.jpg" width="200" alt="编辑备注"/>
</div>

### 🏠 朋友圈与群组
<div align="center">
  <img src="./static/mobile/moments.jpg" width="200" alt="朋友圈"/>
  <img src="./static/mobile/send-moments.jpg" width="200" alt="发布朋友圈"/>
  <img src="./static/mobile/group-list.jpg" width="200" alt="群聊列表"/>
  <img src="./static/mobile/create-group.jpg" width="200" alt="创建群聊"/>
  <img src="./static/mobile/group-details.jpg" width="200" alt="群聊详情"/>
  <img src="./static/mobile/add-members.jpg" width="200" alt="添加成员"/>
</div>

### 👤 个人中心
<div align="center">
  <img src="./static/mobile/mine.jpg" width="200" alt="我的主界面"/>
  <img src="./static/mobile/profile-edit.jpg" width="200" alt="编辑个人资料"/>
  <img src="./static/mobile/qcode.jpg" width="200" alt="二维码功能"/>
</div>

### ⚙️ 系统功能
<div align="center">
  <img src="./static/mobile/settings.jpg" width="200" alt="设置"/>
  <img src="./static/mobile/update.jpg" width="200" alt="更新"/>
  <img src="./static/mobile/feedback.jpg" width="200" alt="反馈"/>
  <img src="./static/mobile/about.jpg" width="200" alt="关于"/>
  <img src="./static/mobile/statement.jpg" width="200" alt="声明"/>
</div>

## 🤝 贡献指南

我们欢迎社区贡献！详情请参阅[贡献指南](CONTRIBUTING.md)。

### 如何贡献

1. **Fork** 本仓库
2. **创建** 特性分支 (`git checkout -b feature/AmazingFeature`)
3. **提交** 更改 (`git commit -m 'Add some AmazingFeature'`)
4. **推送** 到分支 (`git push origin feature/AmazingFeature`)
5. **开启** Pull Request

### 贡献领域

- 🐛 **Bug报告** - 帮助我们识别和修复问题
- 💡 **功能建议** - 建议新功能和改进
- 📝 **文档** - 改进文档和示例
- 🔧 **代码贡献** - 提交代码改进和新功能
- 🧪 **测试** - 帮助测试和质量保证

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

> 📖 **详细法律条款**：请参阅 [LEGAL.md](LEGAL.md) 文件了解完整的法律免责声明和使用要求。

## ⭐ Star历史

[![Star History Chart](https://api.star-history.com/svg?repos=wsrh8888/beaver-server&type=Date)](https://star-history.com/#wsrh8888/beaver-server&Date)

## ☕ 请作者喝杯茶

如果这个项目对你有帮助，欢迎请作者喝杯茶 ☕

<div align="center">
  <img src="./static/sponsor/wechat.jpg" width="200" alt="微信赞助码"/>
  <img src="./static/sponsor/zhifubao.jpg" width="200" alt="支付宝赞助码"/>
</div>

## ⭐ 支持项目

如果这个项目对你有帮助，请在GitHub上给我们一个 ⭐ Star！

---

<div align="center">
  <strong>Made with ❤️ by Beaver IM Team</strong><br>
  <em>企业级即时通讯平台</em>
</div> 