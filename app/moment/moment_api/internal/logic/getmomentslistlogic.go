package logic

import (
	"context"
	"fmt"

	"beaver/app/friend/friend_rpc/types/friend_rpc"
	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/app/user/user_rpc/types/user_rpc"

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
		List:  make([]types.MomentListItem, 0, len(moments)),
	}

	// 获取所有动态UUID，用于批量查询评论和点赞数量
	momentUUIDs := make([]string, 0, len(moments))
	userIds := make(map[string]bool) // 用于去重用户ID
	for _, moment := range moments {
		momentUUIDs = append(momentUUIDs, moment.UUID)
		userIds[moment.UserID] = true
	}

	// 将用户ID转换为数组
	var userIdList []string
	for userId := range userIds {
		userIdList = append(userIdList, userId)
	}

	// 批量查询评论数量
	commentCounts := make(map[string]int64)
	if len(momentUUIDs) > 0 {
		var commentStats []struct {
			MomentID string
			Count    int64
		}
		l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
			Where("moment_id IN (?) AND is_deleted = false", momentUUIDs).
			Select("moment_id, COUNT(*) as count").
			Group("moment_id").
			Scan(&commentStats)

		for _, stat := range commentStats {
			commentCounts[stat.MomentID] = stat.Count
		}
	}

	// 批量查询点赞数量
	likeCounts := make(map[string]int64)
	if len(momentUUIDs) > 0 {
		var likeStats []struct {
			MomentID string
			Count    int64
		}
		l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
			Where("moment_id IN (?) AND is_deleted = false", momentUUIDs).
			Select("moment_id, COUNT(*) as count").
			Group("moment_id").
			Scan(&likeStats)

		for _, stat := range likeStats {
			likeCounts[stat.MomentID] = stat.Count
		}
	}

	// 批量查询当前用户是否点赞
	likedMap := make(map[string]bool)
	if len(momentUUIDs) > 0 {
		var liked []struct {
			MomentID string
		}
		l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
			Where("moment_id IN (?) AND user_id = ? AND is_deleted = false", momentUUIDs, req.UserID).
			Select("moment_id").
			Scan(&liked)
		for _, item := range liked {
			likedMap[item.MomentID] = true
		}
	}

	// 批量查询具体的评论数据（限制每条动态最多返回3条评论）
	commentMap := make(map[string][]moment_models.MomentCommentModel)
	if len(momentUUIDs) > 0 {
		var allComments []moment_models.MomentCommentModel
		// 仅取顶层评论，直接查询所有相关的评论，然后按动态分组
		l.svcCtx.DB.Where("moment_id IN (?) AND is_deleted = false AND parent_id = ''", momentUUIDs).
			Order("moment_id, created_at DESC").
			Find(&allComments)

		// 按动态分组，最多保留3条最新评论
		commentCount := make(map[string]int)
		for _, comment := range allComments {
			if commentCount[comment.MomentID] < 3 {
				commentMap[comment.MomentID] = append(commentMap[comment.MomentID], comment)
				commentCount[comment.MomentID]++
			}
		}
	}

	// 批量查询具体的点赞数据（限制每条动态最多返回10个点赞）
	likeMap := make(map[string][]moment_models.MomentLikeModel)
	if len(momentUUIDs) > 0 {
		var allLikes []moment_models.MomentLikeModel
		// 直接查询所有相关的点赞，然后按动态分组
		l.svcCtx.DB.Where("moment_id IN (?) AND is_deleted = false", momentUUIDs).
			Order("moment_id, created_at DESC").
			Find(&allLikes)

		// 按动态分组，最多保留10个最新点赞
		likeCount := make(map[string]int)
		for _, like := range allLikes {
			if likeCount[like.MomentID] < 10 {
				likeMap[like.MomentID] = append(likeMap[like.MomentID], like)
				likeCount[like.MomentID]++
			}
		}
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

	for _, moment := range moments {
		var files []types.GetMomentsFileInfo
		if moment.Files != nil {
			for _, file := range *moment.Files {
				files = append(files, types.GetMomentsFileInfo{
					FileKey: file.FileKey,
					Type:    uint32(file.Type),
				})
			}
		}

		// 获取用户信息
		userName := ""
		avatar := ""
		if userInfo, exists := userInfoMap[moment.UserID]; exists {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		// 转换评论数据
		comments := convertListComments(commentMap[moment.UUID], userInfoMap)

		// 转换点赞数据
		likes := convertListLikes(likeMap[moment.UUID], userInfoMap)

		// 构建完整响应
		resp.List = append(resp.List, types.MomentListItem{
			Id:           moment.UUID,
			UserID:       moment.UserID,
			Content:      moment.Content,
			Files:        files,
			Comments:     comments,
			Likes:        likes,
			UserName:     userName,
			Avatar:       avatar,
			CommentCount: commentCounts[moment.UUID],
			LikeCount:    likeCounts[moment.UUID],
			IsLiked:      likedMap[moment.UUID],
			CreatedAt:    moment.CreatedAt.String(),
		})
	}

	return resp, nil
}

func convertListComments(comments []moment_models.MomentCommentModel, userInfoMap map[string]*user_rpc.UserInfo) []types.GetMomentsCommentInfo {
	var result []types.GetMomentsCommentInfo
	for _, comment := range comments {
		userName := ""
		avatar := ""
		if userInfo, exists := userInfoMap[comment.UserID]; exists {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		result = append(result, types.GetMomentsCommentInfo{
			Id:               comment.UUID,
			UserID:           comment.UserID,
			UserName:         userName,
			Avatar:           avatar,
			Content:          comment.Content,
			ParentId:         "", // 列表只返回顶层
			ReplyToCommentId: "",
			ReplyToUserName:  "",
			CreatedAt:        comment.CreatedAt.String(),
		})
	}
	return result
}

func convertListLikes(likes []moment_models.MomentLikeModel, userInfoMap map[string]*user_rpc.UserInfo) []types.GetMomentsLikeInfo {
	var result []types.GetMomentsLikeInfo
	for _, like := range likes {
		userName := ""
		avatar := ""
		if userInfo, exists := userInfoMap[like.UserID]; exists {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		result = append(result, types.GetMomentsLikeInfo{
			Id:        like.UUID,
			UserID:    like.UserID,
			CreatedAt: like.CreatedAt.String(),
			UserName:  userName,
			Avatar:    avatar,
		})
	}
	return result
}
