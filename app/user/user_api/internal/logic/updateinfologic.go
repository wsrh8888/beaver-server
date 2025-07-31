package logic

import (
	"context"
	"fmt"

	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"beaver/common/ajax"
	"beaver/common/wsEnum/wsCommandConst"
	"beaver/common/wsEnum/wsTypeConst"

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
	if req.FileName != nil {
		updateFields["file_name"] = *req.FileName
	}
	if req.Abstract != nil {
		updateFields["abstract"] = *req.Abstract
	}
	if req.Gender != nil {
		updateFields["gender"] = *req.Gender
	}

	// 执行更新操作
	if len(updateFields) > 0 {
		err = l.svcCtx.DB.Model(&user).Updates(updateFields).Error
		if err != nil {
			return nil, err
		}
	}

	// 异步更新缓存和通知好友
	go func() {
		// 创建新的context，避免使用请求的context
		ctx := context.Background()
		// 拿到自己的好友列表
		response, err := l.svcCtx.FriendRpc.GetFriendIds(ctx, &friend_rpc.GetFriendIdsRequest{
			UserID: req.UserID,
		})
		if err != nil {
			logx.Errorf("failed to get friend ids: %v", err)
			return
		}
		fmt.Println("转发给好友的列表", response.FriendIds)
		// 通过ws推送给自己的好友
		for _, friendID := range response.FriendIds {
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.USER_PROFILE, wsTypeConst.ProfileChangeNotify, req.UserID, friendID, map[string]interface{}{
				"userId": req.UserID,
			}, "")
		}
	}()

	return &types.UpdateInfoRes{}, nil
}
