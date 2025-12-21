package logic

import (
	"context"
	"fmt"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentLikesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态点赞列表的接口（分页）
func NewGetMomentLikesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentLikesLogic {
	return &GetMomentLikesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentLikesLogic) GetMomentLikes(req *types.GetMomentLikesReq) (resp *types.GetMomentLikesRes, err error) {
	// 分页参数处理
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}
	offset := (page - 1) * limit

	// 获取点赞总数
	var totalCount int64
	if err := l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
		Where("moment_id = ? AND is_deleted = false", req.MomentID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 获取分页点赞数据
	var likes []moment_models.MomentLikeModel
	if err := l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false", req.MomentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&likes).Error; err != nil {
		return nil, err
	}

	// 获取点赞用户ID列表
	userIds := make(map[string]bool)
	for _, like := range likes {
		userIds[like.UserID] = true
	}

	var userIdList []string
	for userId := range userIds {
		userIdList = append(userIdList, userId)
	}

	// 批量获取用户信息
	userInfoMap := make(map[string]*user_rpc.UserInfo)
	if len(userIdList) > 0 {
		userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
			UserIdList: userIdList,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %v", err)
		}
		userInfoMap = userResp.UserInfo
	}

	// 转换点赞数据
	var likeInfos []types.GetMomentLikesInfo
	for _, like := range likes {
		userInfo := userInfoMap[like.UserID]
		userName := ""
		avatar := ""
		if userInfo != nil {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		likeInfos = append(likeInfos, types.GetMomentLikesInfo{
			Id:        like.LikeID,
			UserID:    like.UserID,
			UserName:  userName,
			Avatar:    avatar,
			CreatedAt: like.CreatedAt.String(),
		})
	}

	// 构建响应
	resp = &types.GetMomentLikesRes{
		Count: totalCount,
		List:  likeInfos,
	}

	return resp, nil
}
