package logic

import (
	"context"
	"errors"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteFriendVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除好友验证记录
func NewDeleteFriendVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendVerifyLogic {
	return &DeleteFriendVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFriendVerifyLogic) DeleteFriendVerify(req *types.DeleteFriendVerifyReq) (resp *types.DeleteFriendVerifyRes, err error) {
	// 转换ID
	verifyID, err := strconv.ParseUint(req.VerifyID, 10, 32)
	if err != nil {
		logx.Errorf("无效的好友验证ID: %s", req.VerifyID)
		return nil, errors.New("无效的好友验证ID")
	}

	// 先查询好友验证记录是否存在
	var verify friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Where("id = ?", uint(verifyID)).First(&verify).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("好友验证记录不存在, Id: %s", req.VerifyID)
			return nil, errors.New("好友验证记录不存在")
		}
		logx.Errorf("查询好友验证记录失败: %v", err)
		return nil, err
	}

	// 删除好友验证记录
	err = l.svcCtx.DB.Delete(&verify).Error
	if err != nil {
		logx.Errorf("删除好友验证记录失败: %v", err)
		return nil, err
	}

	logx.Infof("好友验证记录删除成功, Id: %s, SendUserID: %s, RevUserID: %s",
		req.VerifyID, verify.SendUserID, verify.RevUserID)
	return &types.DeleteFriendVerifyRes{}, nil
}
