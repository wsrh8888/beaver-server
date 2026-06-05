package logic

import (
	"context"
	"errors"

	"beaver/app/friend/friend_models"
	"beaver/app/friend/friend_rpc/internal/svc"
	"beaver/app/friend/friend_rpc/types/friend_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const friendVerifyActionDelete int32 = 1 // 删除好友验证记录

type UpdateFriendVerifiesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendVerifiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendVerifiesLogic {
	return &UpdateFriendVerifiesLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdateFriendVerifiesLogic) UpdateFriendVerifies(in *friend_rpc.UpdateFriendVerifiesReq) (*friend_rpc.UpdateFriendVerifiesRes, error) {
	if in.Action != friendVerifyActionDelete {
		return nil, errors.New("不支持的操作类型")
	}

	var ids []uint
	for _, vid := range in.VerifyIds {
		v, err := findFriendVerify(l.svcCtx.DB, vid)
		if err != nil {
			continue
		}
		ids = append(ids, v.Id)
	}
	if len(ids) == 0 {
		return &friend_rpc.UpdateFriendVerifiesRes{}, nil
	}

	if err := l.svcCtx.DB.Where("id IN ?", ids).Delete(&friend_models.FriendVerifyModel{}).Error; err != nil {
		l.Errorf("删除好友验证失败: %v", err)
		return nil, err
	}
	return &friend_rpc.UpdateFriendVerifiesRes{AffectedCount: int64(len(ids))}, nil
}
