package call_models

/*
RTC 实时信令事件说明 (用于 WebSocket 瞬时通知，不进入持久化聊天记录表)：
1. RTC_INVITE  (呼叫): 发起通话，被叫方收到此信号显示振铃界面。
2. RTC_CANCEL  (取消): 发起方在对方接听前主动挂断，被叫方停止振铃。
3. RTC_ACCEPT  (接听): 被叫方接听，发起方收到后进入通话状态。
4. RTC_REJECT  (拒绝): 被叫方拒绝接听，发起方显示对方拒绝并退出。
5. RTC_HANGUP  (挂断): 通话过程中任意一方挂断，所有人退出房间。
*/

const (
	SignalInvite = "RTC_INVITE"   // 呼叫邀请
	SignalCancel = "RTC_CANCEL"   // 取消呼叫
	SignalAccept = "RTC_ACCEPTED" // 接听通话
	SignalReject = "RTC_REJECT"   // 拒绝接听
	SignalHangup = "RTC_HANGUP"   // 通话中挂断
)
