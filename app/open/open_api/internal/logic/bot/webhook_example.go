package bot

// Webhook 使用示例

/*
1. 配置 Webhook

POST /open/webhook/config
{
  "appId": "your_app_id",
  "eventType": "message.receive",
  "targetUrl": "https://your-server.com/webhook",
  "secret": "your_secret_key"  // 可选,不传会自动生成
}

2. 接收 Webhook 事件

你的服务器会收到 POST 请求:

Headers:
  Content-Type: application/json
  X-Beaver-Signature: hmac-sha256-signature
  X-Beaver-Event-Type: message.receive

Body:
{
  "event_id": "evt_1234567890",
  "event_type": "message.receive",
  "timestamp": 1234567890,
  "app_id": "your_app_id",
  "data": {
    "message_id": "msg_abc123",
    "sender_id": "user_123",
    "content": "Hello!",
    "created_at": 1234567890
  }
}

3. 验证签名

import (
  "crypto/hmac"
  "crypto/sha256"
  "encoding/hex"
)

func verifySignature(payload string, signature string, secret string) bool {
  mac := hmac.New(sha256.New, []byte(secret))
  mac.Write([]byte(payload))
  expectedSig := hex.EncodeToString(mac.Sum(nil))
  return hmac.Equal([]byte(signature), []byte(expectedSig))
}

4. 支持的事件类型

- message.receive     收到消息
- group.member.change 群成员变更
- user.status.change  用户状态变更

5. 重试机制

- 默认重试 3 次
- 指数退避: 1s, 2s, 4s, 8s...
- 超时时间可配置(默认 5 秒)

6. 查看推送日志

GET /open/webhook/logs?appId=xxx&eventType=message.receive&page=1&pageSize=20

返回:
{
  "total": 100,
  "list": [
    {
      "id": "1",
      "eventType": "message.receive",
      "payload": "{...}",
      "responseCode": 200,
      "retryCount": 0,
      "status": 1,  // 1成功 2失败
      "createdAt": 1234567890
    }
  ]
}
*/
