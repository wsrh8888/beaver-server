// Code generated by goctl. DO NOT EDIT.
package types

type IBodySendMsg struct {
	ConversationId string `json:"conversationId"`
	Msg            string `json:"msg"`
}

type Msg struct {
	Type    uint   `json:"type"`             //消息类型 1:文本 2:图片 3:视频 4:文件 5、语音 6：语音通话 7：视频通话 8撤回消息 9：回复消息 10：引用消息
	TextMsg string `json:"textMsg,optional"` //文本消息
	ImgMsg  string `json:"imgMsg,optional"`  //图片消息
}

type ProxySendMsgReq struct {
	UserId   string                 `header:"Beaver-User-Id"`
	Command  string                 `json:"command"`
	TargetId string                 `json:"targetId"`
	Type     string                 `json:"type"`
	Body     map[string]interface{} `json:"body"`
}

type ProxySendMsgRes struct {
}

type SendMsgReq struct {
	UserId         string `header:"Beaver-User-Id"`
	ConversationId string `json:"conversationId"`
	Msg            string `json:"msg"`
}

type WsReq struct {
	UserId string `header:"Beaver-User-Id"`
	Token  string `header:"token"`
}

type WsRes struct {
}
