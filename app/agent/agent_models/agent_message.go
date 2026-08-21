package agent_models

import (
	"database/sql/driver"
	"encoding/json"

	"beaver/common/models"
	"beaver/common/models/ctype"
)

// 消息角色
const (
	AgentMsgRoleUser      = "user"      // 用户发送
	AgentMsgRoleAssistant = "assistant" // AI 回复
	AgentMsgRoleSystem    = "system"    // 系统提示（如开场白）
)

// AgentMessage 用户与某 Agent 的聊天记录（含用户消息与 AI 回复）。
// 内容：MsgType + Msg(JSON)；AI 元数据：LLMMeta(JSON 子结构，非独立表)。
// 一用户 × 一 Agent 一条对话流，按 seq 递增；历史在 agent 域，不进 chat。
type AgentMessage struct {
	models.Model
	MessageID  string        `gorm:"size:64;uniqueIndex;not null;comment:消息ID" json:"messageId"`
	AgentID    string        `gorm:"size:64;index:idx_agent_user_seq,priority:1;not null;comment:Agent业务ID" json:"agentId"`
	UserID      string        `gorm:"size:64;index:idx_agent_user_seq,priority:2;index:idx_user_client_msg,priority:1;not null;comment:会话归属用户" json:"userId"`
	Role        string        `gorm:"size:16;index;not null;comment:user用户/assistant AI回复/system系统" json:"role"`
	MsgType     ctype.MsgType `gorm:"not null;comment:消息类型(同chat:文本/图片/markdown等)" json:"msgType"`
	Msg         *ctype.Msg    `gorm:"type:json;comment:消息内容JSON(同chat)" json:"msg"`
	MsgPreview  string        `gorm:"size:200;comment:列表预览" json:"msgPreview"`
	ClientMsgID string        `gorm:"size:64;index:idx_user_client_msg,priority:2;comment:客户端幂等ID" json:"clientMsgId"`
	Status      int8          `gorm:"type:tinyint;not null;default:1;comment:1正常2撤回" json:"status"`
	Seq         int64         `gorm:"not null;default:0;index:idx_agent_user_seq,priority:3;comment:对话内序号从1递增" json:"seq"`
	LLMMeta     *AgentLLMMeta `gorm:"type:json;comment:AI生成元数据JSON(模型/token/耗时等)" json:"llmMeta,omitempty"`
}

func (AgentMessage) TableName() string {
	return "agent_messages"
}

// AgentLLMMeta AI 回复元数据，作为 AgentMessage.LLMMeta 的子 JSON，不是独立表。
type AgentLLMMeta struct {
	Provider         string `json:"provider,omitempty"`         // 供应商：deepseek/openai/anthropic 等
	ModelName        string `json:"modelName,omitempty"`        // 模型名
	PromptTokens     int    `json:"promptTokens,omitempty"`     // 输入 token
	CompletionTokens int    `json:"completionTokens,omitempty"` // 输出 token
	TotalTokens      int    `json:"totalTokens,omitempty"`      // 合计 token
	CachedTokens     int    `json:"cachedTokens,omitempty"`     // 缓存命中 token（若有）
	ReasoningTokens  int    `json:"reasoningTokens,omitempty"`  // 推理 token（若有）
	LatencyMs        int64  `json:"latencyMs,omitempty"`        // 生成耗时 ms
	FinishReason     string `json:"finishReason,omitempty"`     // stop/length/tool_calls/content_filter 等
	TraceID          string `json:"traceId,omitempty"`          // 链路/请求追踪
}

func (c *AgentLLMMeta) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	b, ok := val.([]byte)
	if !ok {
		s, ok := val.(string)
		if !ok {
			return nil
		}
		b = []byte(s)
	}
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, c)
}

func (c AgentLLMMeta) Value() (driver.Value, error) {
	b, err := json.Marshal(c)
	return string(b), err
}
