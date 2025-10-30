package logic

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/friend/friend_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetFriendVerifyDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取好友验证详情
func NewGetFriendVerifyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendVerifyDetailLogic {
	return &GetFriendVerifyDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFriendVerifyDetailLogic) GetFriendVerifyDetail(req *types.GetFriendVerifyDetailReq) (resp *types.GetFriendVerifyDetailRes, err error) {
	// 转换ID
	verifyID, err := strconv.ParseUint(req.VerifyID, 10, 32)
	if err != nil {
		logx.Errorf("无效的好友验证ID: %s", req.VerifyID)
		return nil, errors.New("无效的好友验证ID")
	}

	var verify friend_models.FriendVerifyModel
	err = l.svcCtx.DB.Where("id = ?", uint(verifyID)).
		Preload("SendUserModel").
		Preload("RevUserModel").
		First(&verify).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("好友验证不存在, Id: %s", req.VerifyID)
			return nil, errors.New("好友验证不存在")
		}
		logx.Errorf("查询好友验证详情失败: %v", err)
		return nil, err
	}

	var sendUserName, sendUserFileName, revUserName, revUserFileName string
	if verify.SendUserModel.UserID != "" {
		sendUserName = verify.SendUserModel.NickName
		sendUserFileName = verify.SendUserModel.FileName
	}
	if verify.RevUserModel.UserID != "" {
		revUserName = verify.RevUserModel.NickName
		revUserFileName = verify.RevUserModel.FileName
	}

	return &types.GetFriendVerifyDetailRes{
		FriendVerifyDetailInfo: types.FriendVerifyDetailInfo{
			Id:               fmt.Sprintf("%d", verify.Id),
			SendUserId:       verify.SendUserID,
			SendUserName:     sendUserName,
			SendUserFileName: sendUserFileName,
			RevUserId:        verify.RevUserID,
			RevUserName:      revUserName,
			RevUserFileName:  revUserFileName,
			SendStatus:       int(verify.SendStatus),
			RevStatus:        int(verify.RevStatus),
			Message:          verify.Message,
			CreateTime:       time.Time(verify.CreatedAt).Format(time.RFC3339),
			UpdateTime:       time.Time(verify.UpdatedAt).Format(time.RFC3339),
		},
	}, nil
}
