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
	case "9":
		*m = CallMsgType
	case "10":
		*m = WithdrawMsgType
	case "11":
		*m = ReplyMsgType
	case "12":
		*m = ForwardMsgType
	case "13":
		*m = MarkdownMsgType
	case "14":
		*m = LinkMsgType
	case "15":
		*m = CloudDocMsgType
	default:
		*m = TextMsgType
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
	/**
	 * @description: 音视频通话消息类型
	 */
	CallMsgType
	/**
	 * @description: 撤回消息类型
	 */
	WithdrawMsgType
	/**
	 * @description: 回复消息类型
	 */
	ReplyMsgType
	/**
	 * @description: 转发消息类型（合并转发/消息集合）
	 */
	ForwardMsgType
	/**
	 * @description: Markdown 富文本消息类型
	 */
	MarkdownMsgType
	/**
	 * @description: 链接卡片消息类型
	 */
	LinkMsgType
	/**
	 * @description: 云文档消息类型（会话内分享云文档卡片，对标飞书云文档 document_token）
	 */
	CloudDocMsgType
)

// 云文档类型（对标飞书 folder / docx / sheet / slides 等）
const (
	CloudDocTypeFolder = 0 // 文件夹
	CloudDocTypeDoc    = 1 // 文档
	CloudDocTypeSheet = 2 // 表格
	CloudDocTypeSlide = 3 // 幻灯片
	CloudDocTypeMind  = 4 // 思维笔记
)

type Msg struct {
	Type            MsgType          `json:"type"`                      // 消息类型 1:文本 2:图片 3:视频 4:文件 5:语音 6:表情 7:通知 8:音频文件 9:通话 10:撤回 11:回复 12:转发 13:markdown 14:链接卡片 15:云文档
	TargetMsgID     string           `json:"targetMsgId,omitempty"`     // 目标消息ID (用于对旧消息的指令：撤回、通话状态变更等)
	AtUserIDs       []string         `json:"atUserIds,omitempty"`       // @的用户ID列表，服务端据此触发定向推送；文本中用@昵称占位，前端扫描渲染高亮
	TextMsg         *TextMsg         `json:"textMsg,omitempty"`         // 文本消息
	ImageMsg        *ImageMsg        `json:"imageMsg,omitempty"`        // 图片消息
	VideoMsg        *VideoMsg        `json:"videoMsg,omitempty"`        // 视频消息
	FileMsg         *FileMsg         `json:"fileMsg,omitempty"`         // 文件消息
	VoiceMsg        *VoiceMsg        `json:"voiceMsg,omitempty"`        // 语音消息（移动端录制的短语音）
	EmojiMsg        *EmojiMsg        `json:"emojiMsg,omitempty"`        // 表情消息
	NotificationMsg *NotificationMsg `json:"notificationMsg,omitempty"` // 通知消息（会话内的通知，如：xxx加入了群聊、xxx创建了群等）
	AudioFileMsg    *AudioFileMsg    `json:"audioFileMsg,omitempty"`    // 音频文件消息（用户上传的音频文件）
	CallMsg         *CallMsg         `json:"callMsg,omitempty"`         // 音视频通话
	WithdrawMsg     *WithdrawMsg     `json:"withdrawMsg,omitempty"`     // 撤回消息
	ReplyMsg        *ReplyMsg        `json:"replyMsg,omitempty"`        // 回复消息
	ForwardMsg      *ForwardMsg      `json:"forwardMsg,omitempty"`      // 转发消息（集合）
	MarkdownMsg     *MarkdownMsg     `json:"markdownMsg,omitempty"`     // Markdown 富文本消息
	LinkMsg         *LinkMsg         `json:"linkMsg,omitempty"`         // 链接卡片消息
	CloudDocMsg     *CloudDocMsg     `json:"cloudDocMsg,omitempty"`     // 云文档消息（分享卡片，正文在文档服务）
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
	Content string `json:"content"` // 文本消息内容
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
	FileKey   string `json:"fileKey"`             // 文件ID（对标飞书 file_key / file_token）
	FileName  string `json:"fileName,omitempty"`  // 原始文件名
	Size      int64  `json:"size,omitempty"`      // 文件大小（字节）
	MimeType  string `json:"mimeType,omitempty"`  // MIME 类型
	Extension string `json:"extension,omitempty"` // 扩展名：docx、pptx、xlsx、pdf 等
	OpenMode  int    `json:"openMode,omitempty"`  // 打开方式：1=下载 2=预览 3=在线编辑（Office）
}

// CloudDocMsg 云文档消息（Type: 15，对标飞书会话内分享云文档）
// IM 仅存引用卡片，文档正文与协同状态由独立文档服务维护。
type CloudDocMsg struct {
	DocID    string `json:"docId"`              // 文档唯一标识
	DocType  int    `json:"docType"`            // 文档类型：见 CloudDocType* 常量
	Title    string `json:"title"`              // 标题（会话列表/卡片展示）
	OwnerID  string `json:"ownerId,omitempty"`  // 文档所有者用户ID
	Perm     int    `json:"perm,omitempty"`     // 接收者权限快照：1=可阅读 2=可编辑 3=可管理
	CoverURL string `json:"coverUrl,omitempty"` // 封面/缩略图
	Revision int64  `json:"revision,omitempty"` // 文档版本号（协同冲突提示）
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
	FileKey   string `json:"fileKey"`          // 表情图片文件ID（Emoji.FileName）
	EmojiID   string `json:"emojiId"`          // 表情ID（Emoji.ID，单个表情时使用）
	PackageID string `json:"packageId"`        // 表情包ID（EmojiPackage.ID，表情包分享时使用）
	Width     int64  `json:"width,omitempty"`  // 表情图片宽度（可选）
	Height    int64  `json:"height,omitempty"` // 表情图片高度（可选）
}

// CallMsg 音视频通话消息结构 (用于聊天记录显示)
type CallMsg struct {
	RoomID   string `json:"roomId"`             // 房间ID
	CallType int    `json:"callType"`           // 通话类型: 1-私聊, 2-群聊
	Status   int    `json:"status"`             // 状态: 1-进行中, 2-已结束
	Duration int64  `json:"duration,omitempty"` // 通话时长(秒)
}

// WithdrawMsg 撤回消息结构 (Type: 10)
type WithdrawMsg struct {
	OriginMsgID string `json:"originMsgId"`         // 被撤回的消息ID
	OriginMsg   *Msg   `json:"originMsg,omitempty"` // 被撤回的消息内容快照（用于重新编辑或审计）
}

// ReplyMsg 回复消息结构 (Type: 11)
type ReplyMsg struct {
	OriginMsgID string `json:"originMsgId"`         // 被回复的消息ID
	OriginMsg   *Msg   `json:"originMsg,omitempty"` // 被回复的消息内容快照
	ReplyMsg    *Msg   `json:"replyMsg"`            // 回复的消息主体对象 (可以是文本、图片等)
}

// ForwardMsg 转发消息结构（大厂标准：轻量化卡片 + 延迟加载）
type ForwardMsg struct {
	Title    string `json:"title"`    // 转发的标题，如 "群聊的聊天记录"
	RecordID string `json:"recordId"` // 核心：指向完整详情的聚合ID（存入独立详情表或OSS）
	Count    int    `json:"count"`    // 包含的消息总数
}

// MarkdownMsg Markdown 富文本消息结构（Type: 13）
type MarkdownMsg struct {
	Content string `json:"content"`          // Markdown 正文
	Title   string `json:"title,omitempty"`  // 会话列表预览标题（可选，为空时取 Content 前50字）
}

// LinkMsg 链接卡片消息结构（Type: 14）
type LinkMsg struct {
	URL      string `json:"url"`               // 跳转链接
	Title    string `json:"title"`             // 标题
	Desc     string `json:"desc,omitempty"`    // 摘要描述
	ImageURL string `json:"imageUrl,omitempty"` // 封面图 URL（可选）
}
