package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

type MsgType uint32

const (
	TextMsgType MsgType = iota + 1
	ImageMsgType
	VideoMsgType
	FileMsgType
	VoiceMsgType
	EmojiMsgType
)

type Msg struct {
	Type     MsgType   `json:"type"`               //消息类型 1:文本 2:图片 3:视频 4:文件 5:语音 6:表情
	TextMsg  *TextMsg  `json:"textMsg,omitempty"`  //文本消息
	ImageMsg *ImageMsg `json:"imageMsg,omitempty"` //图片消息
	VideoMsg *VideoMsg `json:"videoMsg,omitempty"` //视频消息
	FileMsg  *FileMsg  `json:"fileMsg,omitempty"`  //文件消息
	VoiceMsg *VoiceMsg `json:"voiceMsg,omitempty"` //语音消息
	EmojiMsg *EmojiMsg `json:"emojiMsg,omitempty"` //表情消息
	ReplyMsg *ReplyMsg `json:"replyMsg,omitempty"` //回复消息
}

/**
 * @description: 取出来的时候的数据
 */
func (c *Msg) Scan(val interface{}) error {
	err := json.Unmarshal(val.([]byte), c)
	if err != nil {
		return err
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
	FileName string `json:"fileName"` //图片文件ID
	Width    int    `json:"width"`    //图片宽度
	Height   int    `json:"height"`   //图片高度
}

type VideoMsg struct {
	FileName string `json:"fileName"` //视频文件ID
	Width    int    `json:"width"`    //视频宽度
	Height   int    `json:"height"`   //视频高度
	Duration int    `json:"duration"` //视频时长（秒）
}

type FileMsg struct {
	FileName string `json:"fileName"` //文件ID，通过fileName可以从文件模块获取完整信息
}

type VoiceMsg struct {
	FileName string `json:"fileName"` //语音文件ID
	Duration int    `json:"duration"` //语音时长（秒）
}

// 表情消息结构
type EmojiMsg struct {
	FileName  string `json:"fileName"`  // 表情图片文件ID（Emoji.FileName）
	EmojiID   uint32 `json:"emojiId"`   // 表情ID（Emoji.ID，单个表情时使用）
	PackageID uint32 `json:"packageId"` // 表情包ID（EmojiPackage.ID，表情包分享时使用）
}

// 回复消息结构
type ReplyMsg struct {
	ReplyToMessageID string `json:"replyToMessageId"` // 被回复的消息ID
	ReplyToContent   string `json:"replyToContent"`   // 被回复的消息内容预览
	ReplyToSender    string `json:"replyToSender"`    // 被回复消息的发送者昵称
}
