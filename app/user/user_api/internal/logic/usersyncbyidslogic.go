package logic

import (
	"beaver/app/user/user_api/internal/svc"
	"beaver/app/user/user_api/internal/types"
	"beaver/app/user/user_models"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/net/context"
)

type UserSyncByIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户数据同步（通过用户ID列表，大厂方式）
func NewUserSyncByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserSyncByIdsLogic {
	return &UserSyncByIdsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserSyncByIdsLogic) UserSyncByIds(req *types.UserSyncByIdsReq) (resp *types.UserSyncByIdsRes, err error) {
	// 设置默认限制
	limit := req.Limit
	if limit <= 0 {
		limit = 100 // 大厂方式默认批量同步更多数据
	}

	// 去重用户ID列表
	uniqueUserIDs := l.removeDuplicates(req.UserIDs)
	if len(uniqueUserIDs) == 0 {
		return &types.UserSyncByIdsRes{
			Users: []types.UserSyncItem{},
		}, nil
	}

	// 直接查询指定用户的最新信息（关键：不检查版本号）
	var users []user_models.UserModel
	err = l.svcCtx.DB.Where("uuid IN ?", uniqueUserIDs).
		Limit(limit).
		Find(&users).Error
	if err != nil {
		l.Errorf("查询指定用户数据失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var userItems []types.UserSyncItem
	for _, user := range users {
		userItems = append(userItems, types.UserSyncItem{
			UserID:   user.UUID,
			Nickname: user.NickName,
			Avatar:   user.Avatar,
			Abstract: user.Abstract,
			Phone:    user.Phone,
			Email:    user.Email,
			Gender:   user.Gender,
			Status:   user.Status,
			Version:  user.Version,
			CreateAt: time.Time(user.CreatedAt).Unix(),
			UpdateAt: time.Time(user.UpdatedAt).Unix(),
		})
	}

	resp = &types.UserSyncByIdsRes{
		Users: userItems,
	}

	l.Infof("通过ID列表同步用户信息完成，请求用户数: %d, 返回用户数: %d", len(uniqueUserIDs), len(userItems))
	return resp, nil
}

// removeDuplicates 去重
func (l *UserSyncByIdsLogic) removeDuplicates(ids []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			result = append(result, id)
		}
	}
	return result
}
