# 建状态脚本说明

所有建状态脚本都在这里，每个文件**内联自己的真实业务 flow**（Go flow），编排（worker pool + 限速）来构造压测前置数据。产出的数据文件放在 `scripts/stress/` 根目录，供所有压测脚本共用。

**结构原则**：Go flow 唯一消费者就是 setup，所以 flow 内联在对应 setup 文件里，不独立成层。k6 侧 flow（JS）在 `utils/k6/` 给压测脚本跨脚本复用，因为 JS 无包系统必须独立。

路径由 `config.AccountsPath()` / `config.FriendsPath()` 统一定位（与 `.env` 同目录），脚本在任意子目录下运行都写到同一份。三个独立 package：`config`（配置/路径）、`httputil`（HTTP 客户端）、`types`（Account/Friendship 类型）。

## 脚本清单

| 脚本 | 作用 | 内联 flow | 产出 | 前置 |
|------|------|-----------|------|------|
| `users/main.go` | 批量注册测试账号 | registerOne（写验证码→调注册接口） | accounts.json | 无 |
| `tokens/main.go` | 批量登录回填 token | loginOne（调登录接口） | accounts.json | 先跑 users |
| `friends/main.go` | 配对互加好友 + 拿会话ID | makeFriends/addFriend/valid/validList/friendInfo | friends.json | 先跑 users + tokens |

## 运行

在 `scripts/stress/setup` 目录下：

```powershell
go run ./users -count 500 -c 10 -rps 5    # 第 1 步：造账号
go run ./tokens -c 10                       # 第 2 步：拿 token
go run ./friends -c 5                        # 第 3 步：造好友关系 + 会话
```

参数说明：
- `-count N`：注册账号数量（users 专用）
- `-start N`：邮箱序号起始值，避免重复注册冲突（users 专用）
- `-c N`：并发 worker 数（friends 涉及双向请求，建议 5-10 不宜过高）
- `-rps N`：每秒提交上限（users 专用，保护服务端不被打爆）

## 数据文件

- `accounts.json`：账号 + token + userId
- `friends.json`：好友对 + conversationId（会话随好友关系自动创建，无独立建会话接口）

## 扩展

新增建状态脚本（比如造群、造朋友圈）：
1. 建 `setup/<名字>/main.go`，`package main`
2. 内联自己的真实业务 flow（HTTP 细节只在这一个文件里），用 `config.Load()` 读配置、`httputil.PostJSON/GetJSON` 发请求、`types.Account` 等数据类型
3. worker pool + 限速编排（照着 `friends/main.go` 的并发模型）
4. 结果写回 `utils.*Path()` 定位的数据文件
5. 若该业务也要压测，同步在 `utils/k6/<服务>.js` 加 JS 版 flow 给 single/scenarios 复用
