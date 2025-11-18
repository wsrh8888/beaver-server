package logic

import (
	"context"
	"fmt"

	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentsListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentsListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentsListLogic {
	return &GetMomentsListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentsListLogic) GetMomentsList(req *types.GetMomentsReq) (resp *types.GetMomentsRes, err error) {
	// 获取好友列表
	friendResp, err := l.svcCtx.FriendRpc.GetFriendIds(l.ctx, &friend_rpc.GetFriendIdsRequest{
		UserID: req.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %v", err)
	}

	// 将自己的ID添加到好友列表中
	friendIds := append(friendResp.FriendIds, req.UserID)

	// 分页参数计算
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	// 获取自己的动态和好友的动态
	var moments []moment_models.MomentModel
	if err := l.svcCtx.DB.Where("user_id IN (?) AND is_deleted = false", friendIds).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&moments).Error; err != nil {
		return nil, err
	}

	// 获取总数
	var count int64
	if err := l.svcCtx.DB.Model(&moment_models.MomentModel{}).
		Where("user_id IN (?) AND is_deleted = false", friendIds).
		Count(&count).Error; err != nil {
		return nil, err
	}

	// 准备响应数据
	resp = &types.GetMomentsRes{
		Count: count,
		List:  make([]types.MomentModel, 0, len(moments)),
	}

	// 获取所有动态ID，用于批量查询评论和点赞数量
	momentIds := make([]uint, 0, len(moments))
	for _, moment := range moments {
		momentIds = append(momentIds, moment.Id)
	}

	// 批量查询评论数量
	commentCounts := make(map[uint]int64)
	if len(momentIds) > 0 {
		var commentStats []struct {
			MomentID uint
			Count    int64
		}
		l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
			Where("moment_id IN (?) AND is_deleted = false", momentIds).
			Select("moment_id, COUNT(*) as count").
			Group("moment_id").
			Scan(&commentStats)

		for _, stat := range commentStats {
			commentCounts[stat.MomentID] = stat.Count
		}
	}

	// 批量查询点赞数量
	likeCounts := make(map[uint]int64)
	if len(momentIds) > 0 {
		var likeStats []struct {
			MomentID uint
			Count    int64
		}
		l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
			Where("moment_id IN (?) AND is_deleted = false", momentIds).
			Select("moment_id, COUNT(*) as count").
			Group("moment_id").
			Scan(&likeStats)

		for _, stat := range likeStats {
			likeCounts[stat.MomentID] = stat.Count
		}
	}

	for _, moment := range moments {
		var files []types.FileInfo
		if moment.Files != nil {
			for _, file := range *moment.Files {
				files = append(files, types.FileInfo{
					FileKey: file.FileKey,
				})
			}
		}

		// 简化响应，只包含基本信息和统计数据
		resp.List = append(resp.List, types.MomentModel{
			Id:        moment.Id,
			UserID:    moment.UserID,
			Content:   moment.Content,
			Files:     files,
			Likes:     make([]types.MomentLikeModel, 0),    // 列表页不显示具体点赞用户
			Comments:  make([]types.MomentCommentModel, 0), // 列表页不显示具体评论
			UserName:  "",                                  // 需要从用户服务获取
			Avatar:    "",                                  // 需要从用户服务获取
			CreatedAt: moment.CreatedAt.String(),
		})
	}

	return resp, nil
}
