流程一：发起通话 (StartCallLogic)
目标： 用户 A 点击拨打，系统完成“占座”并通知 B。

参数校验：调用 UserRpc 检查 TargetId 是否存在；调用 CallRpc.GetUserStatus 检查对方是否忙碌。
创建记录：调用 CallRpc.CreateSession。
生成唯一 RoomID (UUID)。
在 call_sessions 插入记录（Status: 1-呼叫中）。
在 call_participants 插入 A 和 B 的初始记录。
生成令牌：使用 LiveKit SDK 为 A 生成 RoomToken。
发送信令（关键）：
调用 NotificationRpc 或 ChatRpc 往 WebSocket 发送自定义信令：{ "type": "RTC_INVITE", "roomId": "xxx", "callerId": "A", "callType": 1 }。
返回结果：将 RoomID、RoomToken 和 LiveKitUrl 返回给 A。
流程二：接听通话 (GetTokenLogic)
目标： 用户 B 点击“接听”，系统发放准入许可。

权限校验：调用 CallRpc 检查该 RoomID 是否存在，且 req.UserID (B) 是否在该房间的参与者名单中。
生成令牌：使用 LiveKit SDK 为 B 生成 RoomToken。
状态同步：调用 CallRpc.UpdateParticipantStatus，将 B 的状态改为“已接听”。
发送信令：通过 WS 通知 A：{ "type": "RTC_ACCEPTED", "user": "B" }。此时 A 端 App 停止拨号音，双方建立媒体连接。
流程三：挂断/拒绝 (HangupLogic)
目标： 无论谁主动挂断，确保流程闭环。

业务处理：
调用 CallRpc 将 call_sessions 状态改为“已结束”。
记录挂断原因（主动挂断/拒绝）。
信令同步：
通过 WS 通知对方：{ "type": "RTC_HANGUP", "roomId": "xxx" }。对方收到后自动销毁本地 WebRTC 实例。
媒体清理：
（可选）调用 LiveKit Server API 强制销毁该房间（防止死房间占用带宽）。
流程四：Webhook 自动清理与入库 (LiveKitWebhookLogic)
目标： 解决断网、异常挂断时的状态同步。

校验来源：使用 LiveKit SDK 的 auth.TokenVerifier 校验 Webhook 的合法性。
解析事件：
**participant_joined**：如果这是第一个进房的人，记录 call_sessions 的 StartTime。
**room_finished**：通话彻底结束。
数据沉淀（重要）：
从 Webhook 数据中获取 duration（时长）。
调用 CallRpc.FinalizeSession：更新 MySQL 的 EndTime 和 Duration。
写入聊天历史：调用 ChatRpc.SendMessage，往 A 和 B 的聊天框发一条特殊消息：[通话结束] 时长 05:20。



“Cursor，请实现 StartCallLogic.go。逻辑如下：首先从 Header 获取 Beaver-User-Id 作为发起者。通过 CallRpc 检查目标用户是否忙碌。生成一个 UUID 作为 RoomID。调用 CallRpc 的 CreateSession 接口在数据库创建通话记录。使用 github.com/livekit/protocol/auth 包，配合配置文件中的 ApiKey 和 ApiSecret 为发起者生成一个 VideoGrant 权限的 Token。最后，通过 ChatRpc 发送一条 RTC_INVITE 类型的信令给目标用户。请参考 group 服务的 RPC 调用风格。”