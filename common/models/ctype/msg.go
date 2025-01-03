package ctype

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type MsgType uint32

const (
	TextMsgType MsgType = iota + 1
	ImageMsgType
	VideoMsgType
	FileMsgType
	VoiceMsgType
	VoiceCallMsgType
	VideoCallMsgType
	WithdrawMsgType
	ReplyMsgType
	/**
	 * @description: 引用消息
	 */
	QuoteMsgType
	AtMsgType
)

type Msg struct {
	Type         MsgType       `json:"type"`                   //消息类型 1:文本 2:图片 3:视频 4:文件 5、语音 6：语音通话 7：视频通话 8撤回消息 9：回复消息 10：引用消息
	TextMsg      *TextMsg      `json:"textMsg,omitempty"`      //文本消息
	ImageMsg     *ImageMsg     `json:"imageMsg,omitempty"`     //图片消息
	VideoMsg     *VideoMsg     `json:"videoMsg,omitempty"`     //视频消息
	FileMsg      *FileMsg      `json:"fileMsg,omitempty"`      //文件消息
	VoiceMsg     *VoiceMsg     `json:"voiceMsg,omitempty"`     //语音消息
	VoiceCallMsg *VoiceCallMsg `json:"voiceCallMsg,omitempty"` //语音通话消息
	VideoCallMsg *VideoCallMsg `json:"videoCallMsg,omitempty"` //视频通话消息
	WithdrawMsg  *WithdrawMsg  `json:"withdrawMsg,omitempty"`  //撤回消息
	ReplyMsg     *ReplyMsg     `json:"replyMsg,omitempty"`     //回复消息
	QuoteMsg     *QuoteMsg     `json:"quoteMsg,omitempty"`     //引用消息
	AtMsg        *AtMsg        `json:"atMsg,omitempty"`        //@消息 群聊中使用
}

/**
 * @description: 取出来的时候的数据
 */
func (c *Msg) Scan(val interface{}) error {
	err := json.Unmarshal(val.([]byte), c)
	if err != nil {
		return err
	}
	if c.Type == WithdrawMsgType {
		// 如果这个消息是撤回消息，那么就需要特殊处理
		if c.WithdrawMsg != nil {
			c.WithdrawMsg.OriginMsg = nil
		}
	}
	return nil
}

/**
 * @description: 入库的数据
 */
func (c *Msg) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}

type TextMsg struct {
	Content string `json:"content"` //文本消息内容
}

type ImageMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
}

type VideoMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Time  int32  `json:"time"` //视频时长 单位秒
}

type FileMsg struct {
	Title string `json:"title"`
	Src   string `json:"src"`
	Size  int32  `json:"size"` //文件大小 单位字节
	Type  string `json:"type"` //文件类型
}

type VoiceMsg struct {
	Src  string `json:"src"`
	Time int32  `json:"time"` //语音时长 单位秒
}

type VoiceCallMsg struct {
	StartTime time.Time `json:"startTime"` //通话开始时间
	EndTime   time.Time `json:"endTime"`   //通话结束时间
	EndReason int8      `json:"endReason"` //通话结束原因 0: 发起方挂断 1:接收方挂断 3:未打通

}

type VideoCallMsg struct {
	StartTime time.Time `json:"startTime"` //通话开始时间
	EndTime   time.Time `json:"endTime"`   //通话结束时间
	EndReason int8      `json:"endReason"` //通话结束原因 0: 发起方挂断 1:接收方挂断 3:未打通
}

// 撤回消息
type WithdrawMsg struct {
	Content   string `json:"content"`             //撤回提示词
	MsgId     uint   `json:"msgId"`               //撤回的消息Id
	OriginMsg *Msg   `json:"originMsg,omitempty"` //原消息
}

type ReplyMsg struct {
	MsgId         uint      `json:"msgId"`         //回复的消息Id
	Content       string    `json:"content"`       //回复的文本消息
	Msg           *Msg      `json:"msg,omitempty"` //回复的消息
	UserId        uint      `json:"userId"`        //被回复的用户Id
	UserNickName  string    `json:"userNickName"`  //被回复的用户昵称
	OriginMsgDate time.Time `json:"originMsgDate"` //原消息时间
}

type QuoteMsg struct {
	MsgId         uint      `json:"msgId"`         //回复的消息Id
	Content       string    `json:"content"`       //回复的文本消息
	Msg           *Msg      `json:"msg,omitempty"` //回复的消息
	UserId        uint      `json:"userId"`        //被回复的用户Id
	UserNickName  string    `json:"userNickName"`  //被回复的用户昵称
	OriginMsgDate time.Time `json:"originMsgDate"` //原消息时间
}

/**
 * @description: @消息
 */
type AtMsg struct {
	UserId  uint   `json:"userId"`        //被@的用户Id
	Content string `json:"content"`       //回复的文本消息
	Msg     *Msg   `json:"msg,omitempty"` //回复的消息
}
