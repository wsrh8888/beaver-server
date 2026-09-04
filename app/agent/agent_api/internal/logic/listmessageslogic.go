package logic

import (
	"context"
	"errors"
	"strings"

	"beaver/app/agent/agent_api/internal/svc"
	"beaver/app/agent/agent_api/internal/types"
	"beaver/app/agent/agent_rpc/types/agent_rpc"
	beaverlog "beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type ListMessagesLogic struct {
	logger *beaverlog.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMessagesLogic {
	return &ListMessagesLogic{
		logger: beaverlog.New("list_messages", ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListMessagesLogic) ListMessages(req *types.ListMessagesReq) (*types.ListMessagesRes, error) {
	if req == nil {
		return nil, errors.New("request required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("user id required")
	}
	if strings.TrimSpace(req.AgentId) == "" {
		return nil, errors.New("agentId required")
	}

	res, err := l.svcCtx.AgentRpc.ListMessages(l.ctx, &agent_rpc.ListMessagesReq{
		UserId:    req.UserID,
		AgentId:   req.AgentId,
		BeforeSeq: req.BeforeSeq,
		Limit:     req.Limit,
	})
	if err != nil {
		l.logger.Error(model.LogMsg{Text: "获取智能体消息列表失败", Data: map[string]interface{}{"err": err.Error()}})
		return nil, err
	}

	list := make([]types.AgentMessageItem, 0, len(res.List))
	for _, item := range res.List {
		if item == nil {
			continue
		}
		list = append(list, types.AgentMessageItem{
			MessageId:   item.MessageId,
			AgentId:     item.AgentId,
			UserId:      item.UserId,
			Role:        item.Role,
			MsgType:     item.MsgType,
			MsgJson:     item.MsgJson,
			MsgPreview:  item.MsgPreview,
			Seq:         item.Seq,
			LlmMetaJson: item.LlmMetaJson,
			CreatedAt:   item.CreatedAt,
		})
	}

	return &types.ListMessagesRes{List: list}, nil
}
