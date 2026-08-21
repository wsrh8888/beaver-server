package logic

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"beaver/app/agent/agent_models"
	"beaver/app/agent/agent_rpc/internal/svc"
	"beaver/app/agent/agent_rpc/types/agent_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMessagesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessagesLogic {
	return &ListMessagesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMessagesLogic) ListMessages(in *agent_rpc.ListMessagesReq) (*agent_rpc.ListMessagesRes, error) {
	if in == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(in.UserId) == "" {
		return nil, errors.New("user_id required")
	}
	if strings.TrimSpace(in.AgentId) == "" {
		return nil, errors.New("agent_id required")
	}

	limit := int(in.Limit)
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	q := l.svcCtx.DB.WithContext(l.ctx).
		Where("agent_id = ? AND user_id = ? AND status = ?", in.AgentId, in.UserId, 1)
	if in.BeforeSeq > 0 {
		q = q.Where("seq < ?", in.BeforeSeq)
	}

	var rows []agent_models.AgentMessage
	if err := q.Order("seq DESC").Limit(limit).Find(&rows).Error; err != nil {
		l.Errorf("ListMessages query failed: %v", err)
		return nil, err
	}

	// 返回按 seq 升序，方便客户端直接渲染
	list := make([]*agent_rpc.AgentMessageItem, 0, len(rows))
	for i := len(rows) - 1; i >= 0; i-- {
		list = append(list, toMessageItem(&rows[i]))
	}

	return &agent_rpc.ListMessagesRes{List: list}, nil
}

func toMessageItem(m *agent_models.AgentMessage) *agent_rpc.AgentMessageItem {
	item := &agent_rpc.AgentMessageItem{
		MessageId:  m.MessageID,
		AgentId:    m.AgentID,
		UserId:     m.UserID,
		Role:       m.Role,
		MsgType:    uint32(m.MsgType),
		MsgPreview: m.MsgPreview,
		Seq:        m.Seq,
		CreatedAt:  time.Time(m.CreatedAt).UnixMilli(),
	}
	if m.Msg != nil {
		if b, err := json.Marshal(m.Msg); err == nil {
			item.MsgJson = string(b)
		}
	}
	if m.LLMMeta != nil {
		if b, err := json.Marshal(m.LLMMeta); err == nil {
			item.LlmMetaJson = string(b)
		}
	}
	return item
}
