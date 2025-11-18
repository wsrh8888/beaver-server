package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchDeleteFriendVerifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量删除好友验证记录
func NewBatchDeleteFriendVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchDeleteFriendVerifyLogic {
	return &BatchDeleteFriendVerifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchDeleteFriendVerifyLogic) BatchDeleteFriendVerify(req *types.BatchDeleteFriendVerifyReq) (resp *types.BatchDeleteFriendVerifyRes, err error) {
	if len(req.Ids) == 0 {
		return nil, errors.New("删除的好友验证ID列表不能为空")
	}

	// 转换字符串ID为uint切片
	var verifyIDs []uint
	for _, idStr := range req.Ids {
		verifyID, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			logx.Errorf("无效的好友验证ID: %s", idStr)
			return nil, fmt.Errorf("无效的好友验证ID: %s", idStr)
		}
		verifyIDs = append(verifyIDs, uint(verifyID))
	}

	// 先查询要删除的好友验证记录
	var verifies []friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Where("id IN ?", verifyIDs).Find(&verifies).Error
	if err != nil {
		logx.Errorf("查询要删除的好友验证记录失败: %v", err)
		return nil, err
	}

	if len(verifies) == 0 {
		return nil, errors.New("没有找到要删除的好友验证记录")
	}

	// 批量删除好友验证记录
	err = l.svcCtx.DB.Where("id IN ?", verifyIDs).Delete(&friend_models.FriendVerifyModel{}).Error
	if err != nil {
		logx.Errorf("批量删除好友验证记录失败: %v", err)
		return nil, err
	}

	logx.Infof("批量删除好友验证记录完成, 删除数量: %d", len(verifies))
	for _, verify := range verifies {
		logx.Infof("删除好友验证记录 - Id: %d, SendUserID: %s, RevUserID: %s",
			verify.Id, verify.SendUserID, verify.RevUserID)
	}

	return &types.BatchDeleteFriendVerifyRes{}, nil
}
