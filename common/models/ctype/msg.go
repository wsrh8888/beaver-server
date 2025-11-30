package ctype

import (
	"database/sql/driver"
	"encoding/json"
)

type MsgType uint32

// UnmarshalJSON 自定义JSON反序列化，支持字符串和数字格式
func (m *MsgType) UnmarshalJSON(data []byte) error {
	// 尝试作为数字解析
	var num uint32
	if err := json.Unmarshal(data, &num); err == nil {
		*m = MsgType(num)
		return nil
	}

	// 尝试作为字符串解析
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	// 将字符串转换为数字
	switch str {
	case "1":
		*m = TextMsgType
	case "2":
		*m = ImageMsgType
	case "3":
		*m = VideoMsgType
	case "4":
		*m = FileMsgType
	case "5":
		*m = VoiceMsgType
	case "6":
		*m = EmojiMsgType
	case "7":
		*m = NotificationMsgType
	case "8":
		*m = AudioFileMsgType
	default:
		*m = ImageMsgType
	}
	return nil
}

const (
	/***
	 * @description: 文本消息类型
	 */
	TextMsgType MsgType = iota + 1
	/**
	 * @description: 图片消息类型
	 */
	ImageMsgType
	/**
	 * @description: 视频消息类型
	 */
	VideoMsgType
	/**
	 * @description: 文件消息类型
	 */
	FileMsgType
	/**
	 * @description: 语音消息类型（移动端录制的短语音）
	 */
	VoiceMsgType
	/**
	 * @description: 表情消息类型
	 */
	EmojiMsgType
	/**
	 * @description: 通知消息类型（会话内的通知，如：xxx加入了群聊、xxx创建了群等）
	 */
	NotificationMsgType
	/**
	 * @description: 音频文件消息类型（用户上传的音频文件）
	 */
	AudioFileMsgType
)

type Msg struct {
	Type            MsgType          `json:"type"`                      //消息类型 1:文本 2:图片 3:视频 4:文件 5:语音 6:表情 7:通知消息 8:音频文件
	TextMsg         *TextMsg         `json:"textMsg,omitempty"`         //文本消息
	ImageMsg        *ImageMsg        `json:"imageMsg,omitempty"`        //图片消息
	VideoMsg        *VideoMsg        `json:"videoMsg,omitempty"`        //视频消息
	FileMsg         *FileMsg         `json:"fileMsg,omitempty"`         //文件消息
	VoiceMsg        *VoiceMsg        `json:"voiceMsg,omitempty"`        //语音消息（移动端录制的短语音）
	EmojiMsg        *EmojiMsg        `json:"emojiMsg,omitempty"`        //表情消息
	NotificationMsg *NotificationMsg `json:"notificationMsg,omitempty"` //通知消息（会话内的通知，如：xxx加入了群聊、xxx创建了群等）
	AudioFileMsg    *AudioFileMsg    `json:"audioFileMsg,omitempty"`    //音频文件消息（用户上传的音频文件）
	ReplyMsg        *ReplyMsg        `json:"replyMsg,omitempty"`        //回复消息
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

// NotificationMsg 通知消息结构（会话内的通知，如：xxx加入了群聊、xxx创建了群、添加好友成功等）
type NotificationMsg struct {
	Type   int      `json:"type"`   // 通知类型：1=好友欢迎 2=创建群 3=加入群 4=退出群 5=踢出成员 6=转让群主等
	Actors []string `json:"actors"` // 相关用户ID列表
}

type ImageMsg struct {
	FileKey string `json:"fileKey"`          //图片文件ID
	Width   int    `json:"width,omitempty"`  //图片宽度（可选）
	Height  int    `json:"height,omitempty"` //图片高度（可选）
	Size    int64  `json:"size,omitempty"`   //文件大小（字节，可选）
}

type VideoMsg struct {
	FileKey      string `json:"fileKey"`                //视频文件ID
	Width        int    `json:"width,omitempty"`        //视频宽度（可选）
	Height       int    `json:"height,omitempty"`       //视频高度（可选）
	Duration     int    `json:"duration,omitempty"`     //视频时长（秒，可选）
	ThumbnailKey string `json:"thumbnailKey,omitempty"` //视频封面图文件ID（可选）
	Size         int64  `json:"size,omitempty"`         //文件大小（字节，可选）
}

type FileMsg struct {
	FileKey  string `json:"fileKey"`            //文件ID
	FileName string `json:"fileName,omitempty"` //原始文件名（可选，用于显示）
	Size     int64  `json:"size,omitempty"`     //文件大小（字节，可选）
	MimeType string `json:"mimeType,omitempty"` //MIME类型（可选，如 application/pdf）
}

type VoiceMsg struct {
	FileKey  string `json:"fileKey"`            //语音文件ID
	Duration int    `json:"duration,omitempty"` //语音时长（秒，可选）
	Size     int64  `json:"size,omitempty"`     //文件大小（字节，可选）
}

// AudioFileMsg 音频文件消息结构
type AudioFileMsg struct {
	FileKey  string `json:"fileKey"`            //音频文件ID
	FileName string `json:"fileName,omitempty"` //原始文件名（可选，用于显示）
	Duration int    `json:"duration,omitempty"` //音频时长（秒，可选）
	Size     int64  `json:"size,omitempty"`     //文件大小（字节，可选）
}

// 表情消息结构
type EmojiMsg struct {
	FileKey   string `json:"fileKey"`   // 表情图片文件ID（Emoji.FileName）
	EmojiID   uint32 `json:"emojiId"`   // 表情ID（Emoji.ID，单个表情时使用）
	PackageID uint32 `json:"packageId"` // 表情包ID（EmojiPackage.ID，表情包分享时使用）
}

// 回复消息结构
type ReplyMsg struct {
	ReplyToMessageID string `json:"replyToMessageId"` // 被回复的消息ID
	ReplyToContent   string `json:"replyToContent"`   // 被回复的消息内容预览
	ReplyToSender    string `json:"replyToSender"`    // 被回复消息的发送者昵称
}
