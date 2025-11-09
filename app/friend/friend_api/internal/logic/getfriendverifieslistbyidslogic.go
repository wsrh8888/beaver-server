package logic

import (
	"beaver/app/friend/friend_api/internal/svc"
	"beaver/app/friend/friend_api/internal/types"
	"beaver/app/friend/friend_models"
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendVerifiesListByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取好友验证数据（通过UUID）
func NewGetFriendVerifiesListByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifiesListByIdsLogic {
	return &GetFriendVerifiesListByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendVerifiesListByIdsLogic) GetFriendVerifiesListByIds(req *types.GetFriendVerifiesListByIdsReq) (resp *types.GetFriendVerifiesListByIdsRes, err error) {
	if len(req.Uuids) == 0 {
		return &types.GetFriendVerifiesListByIdsRes{
			FriendVerifies: []types.FriendVerifyById{},
		}, nil
	}

	// 查询指定UUID列表中的好友验证信息
	var friendVerifies []friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Where("uuid IN (?)", req.Uuids).Find(&friendVerifies).Error
	if err != nil {
		l.Errorf("查询好友验证信息失败: uuids=%v, error=%v", req.Uuids, err)
		return nil, err
	}

	l.Infof("查询到 %d 个好友验证信息", len(friendVerifies))

	// 转换为响应格式
	var friendVerifiesList []types.FriendVerifyById
	for _, verify := range friendVerifies {
		friendVerifiesList = append(friendVerifiesList, types.FriendVerifyById{
			UUID:       verify.UUID,
			SendUserID: verify.SendUserID,
			RevUserID:  verify.RevUserID,
			SendStatus: int32(verify.SendStatus),
			RevStatus:  int32(verify.RevStatus),
			Message:    verify.Message,
			Source:     verify.Source,
			Version:    verify.Version,
			CreateAt:   time.Time(verify.CreatedAt).UnixMilli(),
			UpdateAt:   time.Time(verify.UpdatedAt).UnixMilli(),
		})
	}

	return &types.GetFriendVerifiesListByIdsRes{
		FriendVerifies: friendVerifiesList,
	}, nil
}
