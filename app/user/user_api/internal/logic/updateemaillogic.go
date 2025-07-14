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

type UpdateEmailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmailLogic {
	return &UpdateEmailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateEmailLogic) UpdateEmail(req *types.UpdateEmailReq) (resp *types.UpdateEmailRes, err error) {
	// 获取要更新的用户信息
	var user user_models.UserModel
	if err := l.svcCtx.DB.Where("uuid = ?", req.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// 验证邮箱验证码
	err = l.verifyEmailCode(req.NewEmail, req.VerifyCode, "update_email")
	if err != nil {
		return nil, err
	}

	// 检查新邮箱是否已被其他用户使用
	var existingUser user_models.UserModel
	if err := l.svcCtx.DB.Where("email = ? AND uuid != ?", req.NewEmail, req.UserID).First(&existingUser).Error; err == nil {
		return nil, fmt.Errorf("email already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 更新用户邮箱
	err = l.svcCtx.DB.Model(&user).Update("email", req.NewEmail).Error
	if err != nil {
		return nil, err
	}

	// 异步更新缓存和通知好友
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
			ajax.SendMessageToWs(l.svcCtx.Config.Etcd, wsCommandConst.USER_PROFILE, wsTypeConst.ProfileChangeNotify, req.UserID, friendID, map[string]interface{}{
				"userId": req.UserID,
			}, "")
		}
	}()

	return &types.UpdateEmailRes{}, nil
}

// 验证邮箱验证码
func (l *UpdateEmailLogic) verifyEmailCode(email, code, codeType string) error {
	// 从Redis获取存储的验证码
	codeKey := fmt.Sprintf("email_code_%s_%s", email, codeType)
	storedCode, err := l.svcCtx.Redis.Get(codeKey).Result()
	if err != nil {
		return fmt.Errorf("验证码已过期或不存在")
	}

	// 验证验证码
	if storedCode != code {
		return fmt.Errorf("验证码错误")
	}

	// 验证成功后删除验证码
	l.svcCtx.Redis.Del(codeKey)

	return nil
}
