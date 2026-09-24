# 造数据脚本说明

所有造数据脚本都在这里，产出的数据文件 `accounts.json` 放在 `scripts/stress/` 根目录，供所有压测脚本共用。

## 脚本清单

| 脚本 | 作用 | 前置 |
|------|------|------|
| `users.go` | 批量注册测试账号（走真实注册链路） | 无 |
| `tokens.go` | 给已有账号批量登录，回填 token | 先跑 users.go |
| `friends.go` | 造好友关系（TODO，压好友/发消息接口前需要） | 先跑 tokens.go |
| `conversations.go` | 造会话（TODO，压发消息接口前需要） | 先跑 friends.go |

## 运行

在 `scripts/stress/seed` 目录下：

```powershell
go run users.go -count 500      # 第 1 步：造账号
go run tokens.go                # 第 2 步：拿 token
```

## 扩展

新增造数据脚本（比如造群、造朋友圈）照着 `users.go` 的结构写即可：
1. `lib.Load()` 读配置
2. 调接口或直连 Redis
3. 结果写回 `../accounts.json` 或新建数据文件
