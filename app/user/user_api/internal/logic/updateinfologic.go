package logic

import (
	"context"
	"fmt"

	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/common/ajax"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateInfoLogic {
	return &UpdateInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateInfoLogic) UpdateInfo(req *types.UpdateInfoReq) (resp *types.UpdateInfoRes, err error) {
	// 获取要更新的用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("uuid = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// 准备更新的字段
	updateFields := make(map[string]interface{})
	if req.Nickname != nil {
		updateFields["nick_name"] = *req.Nickname
	}
	if req.Avatar != nil {
		updateFields["avatar"] = *req.Avatar
	}

	// 执行更新操作
	if len(updateFields) > 0 {
		err = l.svcCtx.DB.Model(&user).Updates(updateFields).Error
		if err != nil {
			return nil, err
		}
	}
	// 异步更新缓存
	defer func() {
		// 拿到自己的好友列表
		response, err := l.svcCtx.FriendRpc.GetFriendIds(l.ctx, &friend_rpc.GetFriendIdsRequest{
			UserID: req.UserID,
		})
		if err != nil {
			logx.Errorf("failed to get friend ids: %v", err)
			return
		}
		fmt.Println("转发给好友的列表", response.FriendIds)
		// 通过ws推送给自己的好友
		for _, friendID := range response.FriendIds {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, "user_update_info", req.UserID, friendID, map[string]interface{}{
				"userId": req.UserID,
			})
		}
	}()

	return nil, nil
}
