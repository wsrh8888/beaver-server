package robot

import (
	"context"
	"errors"
	"fmt"
	"time"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	"beaver/app/open/open_models"
	"beaver/app/open/open_api/internal/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RobotSendMessageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRobotSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RobotSendMessageLogic {
	return &RobotSendMessageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RobotSendMessageLogic) RobotSendMessage(req *types.RobotSendMessageReq, authorization string) (resp *types.RobotSendMessageRes, err error) {
	token, err := utils.ValidateAppAccessToken(l.svcCtx.DB, authorization)
	if err != nil {
		return nil, err
	}
	app, err := utils.LoadAppByID(l.svcCtx.DB, token.AppID)
	if err != nil {
		return nil, err
	}
	if err := utils.RequireAppCapability(app, true, false); err != nil {
		return nil, err
	}

	robot, err := utils.EnsureAppRobot(l.ctx, l.svcCtx.DB, l.svcCtx.UserRpc, app)
	if err != nil {
		return nil, err
	}

	if req.ConversationID == "" || req.Content == "" {
		return nil, errors.New("conversationId 和 content 不能为空")
	}

	if req.IdempotentKey != "" {
		var existing open_models.OpenRobotSendLog
		if err := l.svcCtx.DB.Where("app_id = ? AND idempotent_key = ?", token.AppID, req.IdempotentKey).
			First(&existing).Error; err == nil {
			return &types.RobotSendMessageRes{
				MessageID: existing.MessageID,
				SendTime:  existing.CreatedAt.Unix(),
			}, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("幂等查询失败")
		}
	}

	msgType := req.MsgType
	if msgType == "" {
		msgType = "text"
	}
	replyBody, err := buildRobotChatMsg(msgType, req.Content)
	if err != nil {
		return nil, err
	}
	msg := replyBody
	if req.ReplyMsgID != "" {
		originSnap := &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: "[消息]"},
		}
		getRes, getErr := l.svcCtx.ChatRpc.GetChatMessage(l.ctx, &chat_rpc.GetChatMessageReq{
			ConversationId: req.ConversationID,
			MessageId:      req.ReplyMsgID,
		})
		if getErr == nil && getRes != nil && getRes.Found && getRes.Msg != nil {
			originSnap = getRes.Msg
		}
		msg = &chat_rpc.Msg{
			Type: 11,
			ReplyMsg: &chat_rpc.ReplyMsg{
				OriginMsgId: req.ReplyMsgID,
				OriginMsg:   originSnap,
				ReplyMsg:    replyBody,
			},
		}
	}

	messageID := fmt.Sprintf("robot_%s_%d", token.AppID, time.Now().UnixNano())
	chatRes, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
		UserId:         robot.RobotID,
		ConversationId: req.ConversationID,
		MessageId:      messageID,
		Msg:            msg,
		DeviceId:       "open_robot",
	})
	if err != nil {
		return nil, err
	}

	if req.IdempotentKey != "" {
		_ = l.svcCtx.DB.Create(&open_models.OpenRobotSendLog{
			AppID:          token.AppID,
			IdempotentKey:  req.IdempotentKey,
			MessageID:      chatRes.MessageId,
			ConversationID: req.ConversationID,
		}).Error
	}

	return &types.RobotSendMessageRes{
		MessageID: chatRes.MessageId,
		SendTime:  time.Now().Unix(),
	}, nil
}

func buildRobotChatMsg(msgType, content string) (*chat_rpc.Msg, error) {
	switch msgType {
	case "text", "":
		return &chat_rpc.Msg{
			Type:    1,
			TextMsg: &chat_rpc.TextMsg{Content: content},
		}, nil
	case "markdown":
		return &chat_rpc.Msg{
			Type: 13,
			MarkdownMsg: &chat_rpc.MarkdownMsg{
				Content: content,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported msgType: %s", msgType)
	}
}
