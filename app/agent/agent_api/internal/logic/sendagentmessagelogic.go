package logic

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"unicode/utf8"

	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/app/agent/agent_models"
	"beaver/app/agent/pyagent"
	"beaver/common/models/ctype"
	"beaver/common/sse"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxTextRunes = 5000

type SendAgentMessageLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendAgentMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendAgentMessageLogic {
	return &SendAgentMessageLogic{
		logger: beaverlog.New("send_agent_message", ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SendAgentMessage：
//  1. 本地落用户消息
//  2. 直连 beaver-agent StreamAgentMessage
//  3. done 后本地落 AI 消息
func (l *SendAgentMessageLogic) SendAgentMessage(emit sse.Emitter, req *types.SendAgentMessageReq) {
	if err := validateSendReq(req); err != nil {
		_ = emit("error", map[string]any{"message": err.Error()})
		return
	}

	content := strings.TrimSpace(req.Content)
	userMessageID, seq, err := l.persistUserMessage(req.UserID, req.AgentId, content, req.ClientMsgId)
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "持久化用户消息失败", Data: map[string]interface{}{"err": err.Error()}})
		_ = emit("error", map[string]any{"message": err.Error()})
		return
	}

	if err := emit("accepted", map[string]any{
		"userMessageId": userMessageID,
		"seq":           seq,
		"status":        "accepted",
	}); err != nil {
		return
	}

	modelID, modelJSON, err := l.resolveModelPayload(req)
	if err != nil {
		l.Errorf("resolve model failed: %v", err)
		_ = emit("error", map[string]any{"message": err.Error()})
		return
	}

	stream, err := l.svcCtx.BeaverAgent.StreamAgentMessage(l.ctx, &pyagent.StreamAgentMessageReq{
		UserMessageId: userMessageID,
		AgentId:       req.AgentId,
		UserId:        req.UserID,
		Content:       content,
		DeviceId:      req.DeviceId,
		ModelId:       modelID,
		ModelJson:     modelJSON,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "打开智能体流式回复失败", Data: map[string]interface{}{"err": err.Error()}})
		_ = emit("error", map[string]any{"message": err.Error()})
		return
	}

	var fullReply strings.Builder
	var llmMetaJSON string

	for {
		ev, recvErr := stream.Recv()
		if recvErr == io.EOF {
			break
		}
		if recvErr != nil {
			l.logger.Error(model.LogMsg{Text: "接收智能体流式回复失败", Data: map[string]interface{}{"err": recvErr.Error()}})
			_ = emit("error", map[string]any{"message": recvErr.Error()})
			return
		}
		if ev == nil {
			continue
		}

		switch ev.GetEvent() {
		case "delta":
			fullReply.WriteString(ev.GetContent())
			_ = emit("delta", map[string]any{
				"content": ev.GetContent(),
				"name":    ev.GetName(),
			})
		case "lg":
			_ = emit("lg", map[string]any{
				"name":     ev.GetName(),
				"content":  ev.GetContent(),
				"dataJson": ev.GetDataJson(),
			})
		case "error":
			_ = emit("error", map[string]any{"message": ev.GetErrorMessage()})
			return
		case "done":
			if ev.GetContent() != "" {
				fullReply.Reset()
				fullReply.WriteString(ev.GetContent())
			}
			llmMetaJSON = ev.GetLlmMetaJson()
		}
	}

	reply := strings.TrimSpace(fullReply.String())
	if reply == "" {
		_ = emit("done", map[string]any{
			"userMessageId": userMessageID,
			"status":        "succeeded",
			"message":       "empty reply",
		})
		return
	}

	assistantID, assistantSeq, err := l.persistAssistantMessage(
		req.UserID, req.AgentId, reply, "reply_"+userMessageID, llmMetaJSON,
	)
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "持久化AI回复失败", Data: map[string]interface{}{"err": err.Error()}})
		_ = emit("error", map[string]any{"message": "persist failed: " + err.Error()})
		return
	}

	_ = emit("done", map[string]any{
		"assistantMessageId": assistantID,
		"seq":                assistantSeq,
		"userMessageId":      userMessageID,
		"status":             "succeeded",
		"llmMetaJson":        llmMetaJSON,
	})
}

func (l *SendAgentMessageLogic) persistUserMessage(userID, agentID, content, clientMsgID string) (string, int64, error) {
	if clientMsgID != "" {
		var existing agent_models.AgentMessage
		err := l.svcCtx.DB.WithContext(l.ctx).
			Where("user_id = ? AND client_msg_id = ?", userID, clientMsgID).
			First(&existing).Error
		if err == nil {
			return existing.MessageID, existing.Seq, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, err
		}
	}

	var agent agent_models.Agent
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("agent_id = ? AND user_id = ? AND status = ?", agentID, userID, 1).
		First(&agent).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", 0, errors.New("agent not found or disabled")
	}
	if err != nil {
		return "", 0, err
	}

	msgID := "amsg_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	preview := truncatePreview(content)
	var seq int64

	err = l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		var maxSeq int64
		if err := tx.Model(&agent_models.AgentMessage{}).
			Select("COALESCE(MAX(seq), 0)").
			Where("agent_id = ? AND user_id = ?", agentID, userID).
			Scan(&maxSeq).Error; err != nil {
			return err
		}
		seq = maxSeq + 1
		return tx.Create(&agent_models.AgentMessage{
			MessageID:   msgID,
			AgentID:     agentID,
			UserID:      userID,
			Role:        agent_models.AgentMsgRoleUser,
			MsgType:     ctype.TextMsgType,
			Msg:         &ctype.Msg{Type: ctype.TextMsgType, TextMsg: &ctype.TextMsg{Content: content}},
			MsgPreview:  preview,
			ClientMsgID: clientMsgID,
			Status:      1,
			Seq:         seq,
		}).Error
	})
	return msgID, seq, err
}

func (l *SendAgentMessageLogic) persistAssistantMessage(userID, agentID, content, clientMsgID, llmMetaJSON string) (string, int64, error) {
	if clientMsgID != "" {
		var existing agent_models.AgentMessage
		err := l.svcCtx.DB.WithContext(l.ctx).
			Where("user_id = ? AND client_msg_id = ?", userID, clientMsgID).
			First(&existing).Error
		if err == nil {
			return existing.MessageID, existing.Seq, nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", 0, err
		}
	}

	var llmMeta *agent_models.AgentLLMMeta
	if strings.TrimSpace(llmMetaJSON) != "" {
		var meta agent_models.AgentLLMMeta
		if err := json.Unmarshal([]byte(llmMetaJSON), &meta); err == nil {
			llmMeta = &meta
		}
	}

	msgID := "amsg_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	preview := truncatePreview(content)
	var seq int64

	err := l.svcCtx.DB.WithContext(l.ctx).Transaction(func(tx *gorm.DB) error {
		var maxSeq int64
		if err := tx.Model(&agent_models.AgentMessage{}).
			Select("COALESCE(MAX(seq), 0)").
			Where("agent_id = ? AND user_id = ?", agentID, userID).
			Scan(&maxSeq).Error; err != nil {
			return err
		}
		seq = maxSeq + 1
		return tx.Create(&agent_models.AgentMessage{
			MessageID:   msgID,
			AgentID:     agentID,
			UserID:      userID,
			Role:        agent_models.AgentMsgRoleAssistant,
			MsgType:     ctype.TextMsgType,
			Msg:         &ctype.Msg{Type: ctype.TextMsgType, TextMsg: &ctype.TextMsg{Content: content}},
			MsgPreview:  preview,
			ClientMsgID: clientMsgID,
			Status:      1,
			Seq:         seq,
			LLMMeta:     llmMeta,
		}).Error
	})
	return msgID, seq, err
}

type agentModelPayload struct {
	ModelId   string `json:"modelId"`
	Source    string `json:"source"`
	ModelName string `json:"modelName,omitempty"`
	Endpoint  string `json:"endpoint,omitempty"`
	ApiKey    string `json:"apiKey,omitempty"`
	Provider  string `json:"provider,omitempty"`
}

// resolveModelPayload：按 modelSource/modelId 组装传给 beaver-agent 的 model_json。
// auto / 空：取官方表 sort 最小的上架模型（含 apiKey）；库空则仍下发 auto 由 agent 本地配置兜底。
// official：查官方表拿 modelName/endpoint/apiKey。
// custom：查用户表拿 endpoint+apiKey。
func (l *SendAgentMessageLogic) resolveModelPayload(req *types.SendAgentMessageReq) (string, string, error) {
	modelID := strings.TrimSpace(req.ModelId)
	source := strings.TrimSpace(req.ModelSource)
	if source == "" {
		source = "official"
	}
	if modelID == "" {
		modelID = agent_models.OfficialModelIDAuto
	}

	if source == "official" && modelID == agent_models.OfficialModelIDAuto {
		var row agent_models.AgentOfficialModel
		err := l.svcCtx.DB.WithContext(l.ctx).
			Where("status = ?", 1).
			Order("sort asc, id asc").
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				b, mErr := json.Marshal(agentModelPayload{ModelId: modelID, Source: "official"})
				return modelID, string(b), mErr
			}
			return "", "", err
		}
		b, mErr := json.Marshal(officialPayload(row))
		return row.ModelID, string(b), mErr
	}

	if source == "custom" {
		var row agent_models.AgentUserModel
		err := l.svcCtx.DB.WithContext(l.ctx).
			Where("model_id = ? AND user_id = ? AND status = ?", modelID, strings.TrimSpace(req.UserID), 1).
			First(&row).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return "", "", errors.New("custom model not found")
			}
			return "", "", err
		}
		b, err := json.Marshal(agentModelPayload{
			ModelId:   row.ModelID,
			Source:    "custom",
			ModelName: row.Name,
			Endpoint:  row.Endpoint,
			ApiKey:    row.ApiKey,
		})
		return modelID, string(b), err
	}

	var row agent_models.AgentOfficialModel
	err := l.svcCtx.DB.WithContext(l.ctx).
		Where("model_id = ? AND status = ?", modelID, 1).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("official model not found")
		}
		return "", "", err
	}
	b, err := json.Marshal(officialPayload(row))
	return modelID, string(b), err
}

func officialPayload(row agent_models.AgentOfficialModel) agentModelPayload {
	return agentModelPayload{
		ModelId:   row.ModelID,
		Source:    "official",
		ModelName: row.ModelName,
		Endpoint:  row.Endpoint,
		ApiKey:    row.ApiKey,
		Provider:  row.Provider,
	}
}

func validateSendReq(req *types.SendAgentMessageReq) error {
	if req == nil {
		return errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return errors.New("user id required")
	}
	if strings.TrimSpace(req.AgentId) == "" {
		return errors.New("agentId required")
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return errors.New("content required")
	}
	if utf8.RuneCountInString(content) > maxTextRunes {
		return errors.New("content too long")
	}
	return nil
}

func truncatePreview(content string) string {
	if utf8.RuneCountInString(content) <= 200 {
		return content
	}
	return string([]rune(content)[:200])
}
