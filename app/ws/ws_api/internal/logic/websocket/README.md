# WebSocket 消息协议说明

## 帧类型说明

WS 消息分为两类帧：

| 帧类型   | 适用 Command               | 结构                          |
|----------|---------------------------|-------------------------------|
| 业务帧   | CHAT_MESSAGE / FRIEND_OPERATION / GROUP_OPERATION / CALL | 完整 `command + content.data` 层 |
| 控制帧   | PING / PONG / ACK          | 扁平 JSON，无 content/data 层  |

---

## 一、控制帧（PING / PONG / ACK）

### 结构

```json
{ "command": "PING|PONG|ACK", "messageId": "...", "timestamp": 1712345678000 }
```

字段说明：
- `PING` / `PONG`：携带 `timestamp`，不携带 `messageId`
- `ACK`：携带 `messageId`（对应客户端消息 ID），不携带 `timestamp`

### 1.1 心跳 — 客户端发 PING，服务端回 PONG

```json
// 客户端 → 服务端
{ "command": "PING", "timestamp": 1712345678000 }

// 服务端 → 客户端（echo timestamp）
{ "command": "PONG", "timestamp": 1712345678000 }
```

### 1.2 心跳 — 服务端主动发 PING，客户端回 PONG

```json
// 服务端 → 客户端（定时发送）
{ "command": "PING", "timestamp": 1712345678000 }

// 客户端 → 服务端
{ "command": "PONG", "timestamp": 1712345678000 }
```

### 1.3 ACK — 服务端收到业务消息后立即回执

> ACK 仅表示"服务端已收到"，无失败状态。
> 客户端若在超时内未收到 ACK，应视为网络丢包并重发。

```json
// 服务端 → 客户端（收到任意业务消息后立即发送）
{ "command": "ACK", "messageId": "cli_msg_abc123" }
```

---

## 二、业务帧结构

```json
// 客户端 → 服务端
{
  "command": "<Command>",
  "content": {
    "timestamp": 1712345678000,
    "messageId": "<客户端生成的唯一ID>",
    "data": {
      "type": "<Type>",
      "conversationId": "<会话ID，按需填写>",
      "body": { ... }
    }
  }
}

// 服务端 → 客户端
{
  "code": 0,
  "command": "<Command>",
  "messageId": "<服务端随机8位ID>",
  "serverTime": 1712345678,
  "content": {
    "timestamp": 1712345678000,
    "messageId": "<对应的客户端消息ID，或服务端生成>",
    "data": {
      "type": "<Type>",
      "conversationId": "<会话ID>",
      "body": { ... }
    }
  }
}
```

### Command 枚举

| Command            | 方向            | 说明                          |
|--------------------|-----------------|-------------------------------|
| `CHAT_MESSAGE`     | 双向            | 聊天消息（发送 / 接收 / 同步）|
| `FRIEND_OPERATION` | 服务端 → 客户端 | 好友关系变更推送              |
| `GROUP_OPERATION`  | 服务端 → 客户端 | 群组变更推送                  |
| `USER_PROFILE`     | 服务端 → 客户端 | 用户信息同步                  |
| `NOTIFICATION`     | 服务端 → 客户端 | 通知中心推送                  |
| `EMOJI`            | 服务端 → 客户端 | 表情数据同步                  |
| `CALL`             | 服务端 → 客户端 | 音视频通话信令                |

> `FRIEND_OPERATION` / `GROUP_OPERATION` / `USER_PROFILE` / `NOTIFICATION` / `EMOJI` / `CALL`
> 均由服务端通过 `proxySendMsg` 接口主动推送，**客户端不应主动发送这些 Command**。

---

## 三、CHAT_MESSAGE — 聊天消息

### 3.1 客户端发送私聊文本消息

```json
{
  "command": "CHAT_MESSAGE",
  "content": {
    "timestamp": 1712345678000,
    "messageId": "cli_msg_abc123",
    "data": {
      "type": "private_message_send",
      "conversationId": "conv_u001_u002",
      "body": {
        "conversationId": "conv_u001_u002",
        "msg": {
          "type": 1,
          "textMsg": { "content": "你好！" }
        }
      }
    }
  }
}
```

服务端收到后流程：
1. 立即回 `ACK` `{ "command": "ACK", "messageId": "cli_msg_abc123" }`
2. 存库、推送消息给接收方

### 3.2 客户端发送群聊图片消息

```json
{
  "command": "CHAT_MESSAGE",
  "content": {
    "timestamp": 1712345678000,
    "messageId": "cli_msg_def456",
    "data": {
      "type": "group_message_send",
      "conversationId": "conv_group_g001",
      "body": {
        "conversationId": "conv_group_g001",
        "msg": {
          "type": 2,
          "imageMsg": {
            "fileKey": "img/2024/abc.jpg",
            "width": 1080,
            "height": 720,
            "size": 204800
          }
        }
      }
    }
  }
}
```

### 3.3 服务端推送消息给接收方（`chat_conversation_message_receive`）

```json
{
  "code": 0,
  "command": "CHAT_MESSAGE",
  "messageId": "srv_x9k2mz1a",
  "serverTime": 1712345678,
  "content": {
    "timestamp": 1712345678000,
    "messageId": "srv_push_001",
    "data": {
      "type": "chat_conversation_message_receive",
      "conversationId": "conv_u001_u002",
      "body": {
        "msgId": "srv_msg_uuid_001",
        "conversationId": "conv_u001_u002",
        "senderId": "u001",
        "msg": {
          "type": 1,
          "textMsg": { "content": "你好！" }
        },
        "sendTime": 1712345678000
      }
    }
  }
}
```

---

## 四、FRIEND_OPERATION — 好友关系推送

> 服务端 → 客户端，由其他服务调用 `proxySendMsg` 触发。

### 4.1 好友信息同步

```json
{
  "command": "FRIEND_OPERATION",
  "content": {
    "data": {
      "type": "friend_receive",
      "body": {
        "userId": "u003",
        "nickname": "小明",
        "avatar": "avatar/u003.jpg",
        "remark": ""
      }
    }
  }
}
```

### 4.2 好友验证请求推送

```json
{
  "command": "FRIEND_OPERATION",
  "content": {
    "data": {
      "type": "friend_verify_receive",
      "body": {
        "fromUserId": "u004",
        "message": "我是小红，加个好友吧",
        "verifyId": "verify_uuid_001"
      }
    }
  }
}
```

---

## 五、GROUP_OPERATION — 群组推送

### 5.1 群组信息变更

```json
{
  "command": "GROUP_OPERATION",
  "content": {
    "data": {
      "type": "group_receive",
      "body": {
        "groupId": "g001",
        "name": "Beaver 开发群",
        "avatar": "avatar/group_g001.jpg",
        "memberCount": 42
      }
    }
  }
}
```

### 5.2 群成员变动

```json
{
  "command": "GROUP_OPERATION",
  "content": {
    "data": {
      "type": "group_member_receive",
      "body": { "groupId": "g001", "userId": "u005", "action": "join" }
    }
  }
}
```

### 5.3 入群申请推送（仅群主/管理员）

```json
{
  "command": "GROUP_OPERATION",
  "content": {
    "data": {
      "type": "group_join_request_receive",
      "body": {
        "groupId": "g001",
        "fromUserId": "u006",
        "message": "申请加入",
        "requestId": "req_uuid_001"
      }
    }
  }
}
```

---

## 六、USER_PROFILE — 用户信息同步

### 6.1 用户资料变更

```json
{
  "command": "USER_PROFILE",
  "content": {
    "data": {
      "type": "user_receive",
      "body": { "userId": "u001", "nickname": "Beaver", "avatar": "avatar/u001.jpg" }
    }
  }
}
```

### 6.2 设备被强制下线

```json
{
  "command": "USER_PROFILE",
  "content": {
    "data": {
      "type": "user_kick_receive",
      "body": { "deviceId": "device_desktop_abc", "reason": "账号在其他设备登录" }
    }
  }
}
```

---

## 七、NOTIFICATION — 通知中心

```json
{
  "command": "NOTIFICATION",
  "content": {
    "data": {
      "type": "notification_receive",
      "body": {
        "notificationId": "notif_001",
        "eventType": "friend_request",
        "title": "好友申请",
        "content": "u004 请求添加你为好友",
        "timestamp": 1712345678000
      }
    }
  }
}
```

---

## 八、CALL — 音视频通话信令

```json
{
  "command": "CALL",
  "content": {
    "data": {
      "type": "call_receive",
      "body": {
        "roomId": "room_livekit_abc",
        "callType": 1,
        "fromUserId": "u001",
        "action": "invite"
      }
    }
  }
}
```

---

## 消息体 msg.type 枚举（CHAT_MESSAGE 专用）

| type | 说明       | body 字段         |
|------|------------|-------------------|
| 1    | 文本       | `textMsg.content` |
| 2    | 图片       | `imageMsg`        |
| 3    | 视频       | `videoMsg`        |
| 4    | 文件       | `fileMsg`         |
| 5    | 语音       | `voiceMsg`        |
| 6    | 表情       | `emojiMsg`        |
| 7    | 通知消息   | `notificationMsg` |
| 8    | 音频文件   | `audioFileMsg`    |
| 9    | 音视频通话 | `callMsg`         |
| 10   | 撤回消息   | `withdrawMsg`     |
| 11   | 回复消息   | `replyMsg`        |
| 12   | 转发消息   | `forwardMsg`      |
