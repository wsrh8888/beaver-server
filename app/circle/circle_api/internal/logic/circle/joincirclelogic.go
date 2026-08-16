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

	"github.com/zeromicro/go-zero/core/logx"
)

type JoinCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewJoinCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *JoinCircleLogic {
	return &JoinCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *JoinCircleLogic) JoinCircle(req *types.JoinCircleReq) (resp *types.JoinCircleRes, err error) {
	var circle circle_models.CircleModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", req.CircleID).First(&circle).Error; err != nil {
		return nil, fmt.Errorf("圈子不存在")
	}

	// 已经是成员
	var existing circle_models.CircleMemberModel
	if l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&existing).Error == nil {
		return &types.JoinCircleRes{Status: 1}, nil
	}

	// 审批加入：创建申请记录
	if circle.JoinType == 1 {
		joinReq := circle_models.CircleJoinRequestModel{
			CircleID: req.CircleID,
			UserID:   req.UserID,
			Status:   0,
			Reason:   req.Reason,
		}
		if err = l.svcCtx.DB.Create(&joinReq).Error; err != nil {
			return nil, fmt.Errorf("提交申请失败: %v", err)
		}
		return &types.JoinCircleRes{Status: 0}, nil
	}

	// 自由加入
	member := circle_models.CircleMemberModel{
		CircleID: req.CircleID,
		UserID:   req.UserID,
		Role:     3,
	}
	if err = l.svcCtx.DB.Create(&member).Error; err != nil {
		return nil, fmt.Errorf("加入圈子失败: %v", err)
	}

	// 更新圈子版本
	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Update("version", circleVersion)

	conversationID := fmt.Sprintf("circle_%s", req.CircleID)

	go func() {
		ctx := context.Background()
		l.svcCtx.ChatRpc.InitializeConversation(ctx, &chat_rpc.InitializeConversationReq{
			ConversationId: conversationID,
			Type:           3,
			UserIds:        []string{req.UserID},
		})

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
								"circleId": req.CircleID,
							},
						},
					},
				},
			},
			"conversationId": conversationID,
		}
		if err := l.svcCtx.RocketMQ.SendMessage(ctx, mqwsconst.MqTopicWs, payload); err != nil {
			logx.Errorf("推送圈子资料同步失败: circleID=%s, err=%v", req.CircleID, err)
		}
	}()

	return &types.JoinCircleRes{Status: 1}, nil
}
