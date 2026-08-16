package circle

import (
	"context"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	mqwsconst "beaver/common/const/mqwsconst"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"
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
		Version:     circleVersion,
	}
	if err = l.svcCtx.DB.Create(&circle).Error; err != nil {
		return nil, fmt.Errorf("创建圈子失败: %v", err)
	}

	// 创建者自动成为圈主
	member := circle_models.CircleMemberModel{
		CircleID: circleID,
		UserID:   req.UserID,
		Role:     1,
	}
	if err = l.svcCtx.DB.Create(&member).Error; err != nil {
		return nil, fmt.Errorf("创建圈主成员失败: %v", err)
	}

	conversationID := fmt.Sprintf("circle_%s", circleID)

	go func() {
		ctx := context.Background()

		_, initErr := l.svcCtx.ChatRpc.InitializeConversation(ctx, &chat_rpc.InitializeConversationReq{
			ConversationId: conversationID,
			Type:           3,
			UserIds:        []string{req.UserID},
		})
		if initErr != nil {
			l.logger.Error(model.LogMsg{
				Text: "初始化圈子会话失败",
				Data: map[string]interface{}{"circleId": circleID, "err": initErr.Error()},
			})
			return
		}

		_, notifyErr := l.svcCtx.ChatRpc.SendNotificationMessage(ctx, &chat_rpc.SendNotificationMessageReq{
			ConversationId: conversationID,
			MessageType:    2,
			Content:        fmt.Sprintf("%s 创建了圈子", req.UserID),
			RelatedUserId:  req.UserID,
		})
		if notifyErr != nil {
			logx.Errorf("发送圈子创建通知失败: circleID=%s, err=%v", circleID, notifyErr)
		}

		// 推送圈子资料变更，客户端同步本地 circles 表后会话列表才能 join 名称/头像
		payload := map[string]interface{}{
			"command":  wsCommandConst.CIRCLE_OPERATION,
			"type":     wsTypeConst.CircleReceive,
			"senderId": req.UserID,
			"targetId": req.UserID,
			"body": map[string]interface{}{
				"tables": []map[string]interface{}{
					{
						"table": "circles",
						"data": []map[string]interface{}{
							{
								"version":  circleVersion,
								"circleId": circleID,
							},
						},
					},
				},
			},
			"conversationId": conversationID,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logx.Errorf("推送圈子资料同步失败: circleID=%s, err=%v", circleID, err)
		}
	}()

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
