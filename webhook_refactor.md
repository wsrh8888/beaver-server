# Webhook 优化方案 (大厂实践版)

针对您在 `common/webhook/callback_sender.go` 中提出的问题，目前的实现虽然能通，但确实存在一些隐患。以下是大厂在处理这类需求时的通用做法和优化建议。

## 1. 核心改进点：大厂是怎么做的？

### A. 架构解耦：从 "同步/直连" 到 "异步队列"
目前代码中使用 `go sendCallbackAsync(...)` 开启协程。
- **问题**: 如果服务器重启或崩溃，内存中待发送的回调会全部丢失。而且无法限制并发量，如果回调目标很多，会瞬间瞬间消耗大量资源。
- **大厂做法**: 
    1. **生产者 (Producer)**: 业务逻辑触发事件，将消息发到消息队列（如 Kafka, RocketMQ）或存储到数据库的 `webhook_tasks` 表中。
    2. **消费者 (Consumer)**: 专门的服务监听队列，取消息并执行 HTTP 发送。
    3. **重试机制**: 利用消息队列的延迟重试（Retry Queue）功能，而不是在代码里 `time.Sleep`。

### B. 安全性 (Security)
目前的签名已经不错了，但通常还会增加：
- **时间戳校验**: 在 Header 中增加 `X-Webhook-Timestamp`。接收方校验该时间与当前时间的差距（通常允许 5 分钟误差），防止**重放攻击**。
- **IP 白名单**: 大厂通常会提供固定的出口 IP，让客户配置白名单。

### C. 可观测性 (Observability)
- **详细记录**: 除了状态，还要记录请求体、响应头、响应体、耗时等。
- **管理后台**: 用户可以在后台查看每一个回调的详情，并手动触发 "重新发送"。

---

## 2. 代码重构建议：引入到哪？

### 推荐方案：注入到 `ServiceContext`
不建议在业务 Logic 中直接调用 `webhook.SendCallback`。应该将其封装为 `ServiceContext` 中的一个组件。

#### 修改步骤：

1. **重构 `CallbackSender`**:
   将 `common/webhook/callback_sender.go` 中的函数封装成一个 `CallbackSender` 接口或结构体。

2. **在 `ServiceContext` 中持有**:
   ```go
   type ServiceContext struct {
       Config         config.Config
       DB             *gorm.DB
       // ... 其他 RPC
       WebhookSender  *webhook.CallbackSender // 注入这里
   }
   ```

3. **在 `NewServiceContext` 中初始化**:
   这样 `CallbackSender` 就可以获取到 `Config` 中的超时配置、重试次数，以及用于记录日志的 `DB` 实例。

---

## 3. 具体修复建议

### 针对 `callback_sender.go` 的代码修复：

1. **Header 修复**: 
   ```go
   req.Header.Set("X-Webhook-Event-ID", data.EventID) // 之前是空的
   req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", data.Timestamp))
   ```

2. **注入 DB 实例**:
   不要在函数内部用 `coregorm.GetDB()`。大厂实践中，所有的数据库操作都应该通过上下文传递的 DB 实例进行，这样可以统一事务管理和更好的测试。

3. **限制并发**:
   如果不使用 MQ，建议引入一个 `Worker Pool`（如 `github.com/gammazero/workerpool`），防止瞬间产生数万个 Goroutine。
