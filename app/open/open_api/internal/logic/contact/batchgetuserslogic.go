package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量获取用户
func NewBatchGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUsersLogic {
	return &BatchGetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BatchGetUsersLogic) BatchGetUsers(req *types.BatchGetUsersReq) (resp *types.BatchGetUsersRes, err error) {
	// 1. 参数校验
	if len(req.UserIDs) == 0 {
		return nil, errors.New("用户 ID 列表不能为空")
	}
	if len(req.UserIDs) > 50 {
		return nil, errors.New("最多支持 50 个用户 ID")
	}

	// 2. 查询用户信息
	var users []user_models.UserModel
	if err := l.svcCtx.DB.Where("user_id IN ?", req.UserIDs).Find(&users).Error; err != nil {
		return nil, errors.New("查询用户失败")
	}

	// 3. 转换为响应格式
	var userDetails []types.UserDetailInfo
	for _, user := range users {
		userDetails = append(userDetails, types.UserDetailInfo{
			UserID:   user.UserID,
			Nickname: user.NickName,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
			Email:    user.Email,
		})
	}

	return &types.BatchGetUsersRes{
		Users: userDetails,
	}, nil
}
