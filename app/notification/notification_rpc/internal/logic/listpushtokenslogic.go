package logic

import (
	"context"
	"errors"

	"beaver/app/notification/notification_models"
	"beaver/app/notification/notification_rpc/internal/svc"
	"beaver/app/notification/notification_rpc/types/notification_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPushTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPushTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPushTokensLogic {
	return &ListPushTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPushTokensLogic) ListPushTokens(in *notification_rpc.ListPushTokensReq) (*notification_rpc.ListPushTokensRes, error) {
	if in.UserId == "" {
		return nil, errors.New("userId 不能为空")
	}

	var rows []notification_models.PushRegistrationModel
	if err := l.svcCtx.DB.Where("user_id = ? AND enabled = ?", in.UserId, true).Find(&rows).Error; err != nil {
		l.Errorf("查询 Push Token 失败: userId=%s, err=%v", in.UserId, err)
		return nil, errors.New("查询 Push Token 失败")
	}

	tokens := make([]*notification_rpc.PushTokenInfo, 0, len(rows))
	for _, row := range rows {
		tokens = append(tokens, &notification_rpc.PushTokenInfo{
			DeviceId:     row.DeviceID,
			PushToken:    row.PushToken,
			PushPlatform: row.PushPlatform,
		})
	}
	return &notification_rpc.ListPushTokensRes{Tokens: tokens}, nil
}
