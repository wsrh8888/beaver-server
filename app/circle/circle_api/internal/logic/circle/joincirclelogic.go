package circle

import (
	"context"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

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
	memberVersion := l.svcCtx.VersionGen.GetNextVersion("circle_members", "circle_id", req.CircleID)
	member := circle_models.CircleMemberModel{
		CircleID: req.CircleID,
		UserID:   req.UserID,
		Role:     3,
		Version:  memberVersion,
	}
	if err = l.svcCtx.DB.Create(&member).Error; err != nil {
		return nil, fmt.Errorf("加入圈子失败: %v", err)
	}

	// 更新圈子版本和成员数
	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		Updates(map[string]interface{}{
			"member_count": circle.MemberCount + 1,
			"version":      circleVersion,
		})

	// 初始化该用户的圈子会话
	conversationID := fmt.Sprintf("circle_%s", req.CircleID)
	l.svcCtx.ChatRpc.InitializeConversation(l.ctx, &chat_rpc.InitializeConversationReq{
		ConversationId: conversationID,
		Type:           3,
		UserIds:        []string{req.UserID},
	})

	return &types.JoinCircleRes{Status: 1}, nil
}
