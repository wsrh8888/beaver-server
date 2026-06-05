package logic

import (
	"context"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateFriendBlocksLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendBlocksLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendBlocksLogic {
	return &UpdateFriendBlocksLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateFriendBlocksLogic) UpdateFriendBlocks(in *friend_rpc.UpdateFriendBlocksReq) (*friend_rpc.UpdateFriendBlocksRes, error) {
	if in.Action != 1 {
		return nil, status.Error(codes.InvalidArgument, "无效的操作类型")
	}
	if len(in.BlockIds) == 0 {
		return &friend_rpc.UpdateFriendBlocksRes{}, nil
	}

	result := l.svcCtx.DB.Where("block_id IN ?", in.BlockIds).Delete(&friend_models.FriendBlockModel{})
	if result.Error != nil {
		l.Errorf("解除黑名单失败: %v", result.Error)
		return nil, status.Error(codes.Internal, "操作失败")
	}
	if result.RowsAffected == 0 {
		result = l.svcCtx.DB.Where("id IN ?", in.BlockIds).Delete(&friend_models.FriendBlockModel{})
		if result.Error != nil {
			l.Errorf("解除黑名单失败: %v", result.Error)
			return nil, status.Error(codes.Internal, "操作失败")
		}
	}

	return &friend_rpc.UpdateFriendBlocksRes{AffectedCount: result.RowsAffected}, nil
}
