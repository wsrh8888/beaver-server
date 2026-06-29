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

type HandleJoinRequestLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHandleJoinRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleJoinRequestLogic {
	return &HandleJoinRequestLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HandleJoinRequestLogic) HandleJoinRequest(req *types.HandleJoinRequestReq) (resp *types.HandleJoinRequestRes, err error) {
	// 权限校验
	var joinReq circle_models.CircleJoinRequestModel
	if err = l.svcCtx.DB.Where("id = ? AND status = 0", req.RequestID).First(&joinReq).Error; err != nil {
		return nil, fmt.Errorf("申请记录不存在或已处理")
	}

	var operator circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", joinReq.CircleID, req.UserID).First(&operator).Error; err != nil {
		return nil, fmt.Errorf("无权限")
	}
	if operator.Role > 2 {
		return nil, fmt.Errorf("仅圈主和管理员可处理申请")
	}

	// 更新申请状态
	if err = l.svcCtx.DB.Model(&joinReq).Update("status", req.Status).Error; err != nil {
		return nil, fmt.Errorf("处理申请失败: %v", err)
	}

	// 通过则加入成员
	if req.Status == 1 {
		memberVersion := l.svcCtx.VersionGen.GetNextVersion("circle_members", "circle_id", joinReq.CircleID)
		member := circle_models.CircleMemberModel{
			CircleID: joinReq.CircleID,
			UserID:   joinReq.UserID,
			Role:     3,
			Version:  memberVersion,
		}
		l.svcCtx.DB.Create(&member)

		circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", joinReq.CircleID)
		l.svcCtx.DB.Model(&circle_models.CircleModel{}).
			Where("circle_id = ?", joinReq.CircleID).
			Updates(map[string]interface{}{
				"member_count": l.svcCtx.DB.Raw("member_count + 1"),
				"version":      circleVersion,
			})

		// 初始化圈子会话
		conversationID := fmt.Sprintf("circle_%s", joinReq.CircleID)
		l.svcCtx.ChatRpc.InitializeConversation(l.ctx, &chat_rpc.InitializeConversationReq{
			ConversationId: conversationID,
			Type:           3,
			UserIds:        []string{joinReq.UserID},
		})
	}

	return &types.HandleJoinRequestRes{}, nil
}
