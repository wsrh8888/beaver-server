package circle

import (
	"context"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/utils/logger"
	"beaver/utils/logger/model"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logger.Logger
}

func NewCreateCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCircleLogic {
	return &CreateCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: logger.New("create_circle"),
	}
}

func (l *CreateCircleLogic) CreateCircle(req *types.CreateCircleReq) (resp *types.CreateCircleRes, err error) {
	circleID := uuid.New().String()
	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", circleID)

	circle := circle_models.CircleModel{
		CircleID:    circleID,
		Name:        req.Name,
		Description: req.Description,
		Avatar:      req.Avatar,
		CreatorID:   req.UserID,
		JoinType:    req.JoinType,
		MemberCount: 1,
		Version:     circleVersion,
	}
	if err = l.svcCtx.DB.Create(&circle).Error; err != nil {
		return nil, fmt.Errorf("创建圈子失败: %v", err)
	}

	// 创建者自动成为圈主
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("circle_members", "circle_id", circleID)
	member := circle_models.CircleMemberModel{
		CircleID: circleID,
		UserID:   req.UserID,
		Role:     1,
		Version:  memberVersion,
	}
	if err = l.svcCtx.DB.Create(&member).Error; err != nil {
		return nil, fmt.Errorf("创建圈主成员失败: %v", err)
	}

	// 在会话系统里为创建者初始化圈子会话（ConversationType=3）
	conversationID := fmt.Sprintf("circle_%s", circleID)
	_, err = l.svcCtx.ChatRpc.InitializeConversation(l.ctx, &chat_rpc.InitializeConversationReq{
		ConversationId: conversationID,
		Type:           3,
		UserIds:        []string{req.UserID},
	})
	if err != nil {
		l.logger.Error(model.LogMsg{
			Text: "初始化圈子会话失败",
			Data: map[string]interface{}{"circleId": circleID, "err": err.Error()},
		})
	}

	l.logger.Info(model.LogMsg{
		Text: "圈子创建成功",
		Data: map[string]interface{}{"circleId": circleID, "userId": req.UserID},
	})

	return &types.CreateCircleRes{
		CircleID:    circleID,
		Name:        circle.Name,
		Description: circle.Description,
		Avatar:      circle.Avatar,
		JoinType:    circle.JoinType,
		CreatorID:   circle.CreatorID,
		CreatedAt:   circle.CreatedAt.String(),
	}, nil
}
