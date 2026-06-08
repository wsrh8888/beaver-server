package model

// LogMsg 结构化日志内容
type LogMsg struct {
	Text string      `json:"text"`
	Data interface{} `json:"data,omitempty"`
}
