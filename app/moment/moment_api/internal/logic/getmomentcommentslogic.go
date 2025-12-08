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

type GetMomentCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态评论列表的接口（分页）
func NewGetMomentCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentCommentsLogic {
	return &GetMomentCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentCommentsLogic) GetMomentComments(req *types.GetMomentCommentsReq) (resp *types.GetMomentCommentsRes, err error) {
	// 分页参数处理
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	// 如果带 parentId，查询该顶层评论下的子回复（专用二级评论接口）
	if req.ParentId != "" {
		var parent moment_models.MomentCommentModel
		if err := l.svcCtx.DB.Where("comment_id = ? AND is_deleted = false", req.ParentId).First(&parent).Error; err != nil {
			return nil, err
		}
		if parent.MomentID != req.MomentID {
			return nil, fmt.Errorf("parent comment not belong to moment")
		}
		if parent.ParentID != "" {
			return nil, fmt.Errorf("only support two-level comments")
		}

		// 子评论总数
		var totalCount int64
		if err := l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
			Where("moment_id = ? AND is_deleted = false AND parent_id = ?", req.MomentID, req.ParentId).
			Count(&totalCount).Error; err != nil {
			return nil, err
		}

		// 分页子评论
		var children []moment_models.MomentCommentModel
		if err := l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false AND parent_id = ?", req.MomentID, req.ParentId).
			Order("created_at DESC").
			Offset(offset).
			Limit(limit).
			Find(&children).Error; err != nil {
			return nil, err
		}

		// 收集用户ID（评论用户 + 被回复用户）
		userIds := make(map[string]bool)
		var replyIds []string
		for _, c := range children {
			userIds[c.UserID] = true
			if c.ReplyToCommentID != "" {
				replyIds = append(replyIds, c.ReplyToCommentID)
			}
		}

		replyTargetMap := make(map[string]moment_models.MomentCommentModel)
		if len(replyIds) > 0 {
			var replyTargets []moment_models.MomentCommentModel
			l.svcCtx.DB.Where("comment_id IN (?)", replyIds).Find(&replyTargets)
			for _, rt := range replyTargets {
				replyTargetMap[rt.CommentID] = rt
				userIds[rt.UserID] = true
			}
		}

		var userIdList []string
		for uid := range userIds {
			userIdList = append(userIdList, uid)
		}

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

		var commentInfos []types.GetMomentCommentsInfo
		for _, c := range children {
			userInfo := userInfoMap[c.UserID]
			userName := ""
			avatar := ""
			if userInfo != nil {
				userName = userInfo.NickName
				avatar = userInfo.Avatar
			}

			replyName := ""
			if c.ReplyToCommentID != "" {
				if target, ok := replyTargetMap[c.ReplyToCommentID]; ok {
					if u := userInfoMap[target.UserID]; u != nil {
						replyName = u.NickName
					}
				}
			}

			commentInfos = append(commentInfos, types.GetMomentCommentsInfo{
				Id:               c.CommentID,
				UserID:           c.UserID,
				UserName:         userName,
				Avatar:           avatar,
				Content:          c.Content,
				ParentId:         c.ParentID,
				ReplyToCommentId: c.ReplyToCommentID,
				ReplyToUserName:  replyName,
				ChildCount:       0,
				Children:         nil,
				CreatedAt:        c.CreatedAt.String(),
			})
		}

		resp = &types.GetMomentCommentsRes{
			Count: totalCount,
			List:  commentInfos,
		}
		return resp, nil
	}

	// 未带 parentId，默认查询顶层评论（返回 childCount + 预览 children，预览最多3条，完整二级由前端传 parentId 再拉）
	// 顶层条件：兼容空字符串、NULL、纯空格
	baseTopWhere := "moment_id = ? AND is_deleted = false AND (parent_id IS NULL OR TRIM(parent_id) = '')"
	var totalCount int64
	if err := l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
		Where(baseTopWhere, req.MomentID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	var topComments []moment_models.MomentCommentModel
	if err := l.svcCtx.DB.Where(baseTopWhere, req.MomentID).
		Order("created_at DESC").
		Order("comment_id DESC").
		Offset(offset).
		Limit(limit).
		Find(&topComments).Error; err != nil {
		return nil, err
	}

	// 统计子评论数量
	childCountMap := make(map[string]int64)
	childMap := make(map[string][]moment_models.MomentCommentModel)
	const childPreviewLimit = 3
	if len(topComments) > 0 {
		var childStats []struct {
			ParentID string
			Count    int64
		}
		var parentIDs []string
		for _, c := range topComments {
			parentIDs = append(parentIDs, c.CommentID)
		}
		l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
			Where("parent_id IN (?) AND is_deleted = false", parentIDs).
			Select("parent_id, COUNT(*) as count").
			Group("parent_id").
			Scan(&childStats)
		for _, stat := range childStats {
			childCountMap[stat.ParentID] = stat.Count
		}

		// 取预览 children
		var children []moment_models.MomentCommentModel
		l.svcCtx.DB.Where("parent_id IN (?) AND is_deleted = false", parentIDs).
			Order("created_at DESC").
			Find(&children)
		countPerParent := make(map[string]int)
		for _, ch := range children {
			if countPerParent[ch.ParentID] < childPreviewLimit {
				childMap[ch.ParentID] = append(childMap[ch.ParentID], ch)
				countPerParent[ch.ParentID]++
			}
		}
	}

	// 获取评论用户ID列表 + 预览子评论用户ID + 被回复用户ID
	userIds := make(map[string]bool)
	replyIds := make(map[string]bool)
	for _, comment := range topComments {
		userIds[comment.UserID] = true
	}
	for _, children := range childMap {
		for _, ch := range children {
			userIds[ch.UserID] = true
			if ch.ReplyToCommentID != "" {
				replyIds[ch.ReplyToCommentID] = true
			}
		}
	}

	replyTargetMap := make(map[string]moment_models.MomentCommentModel)
	if len(replyIds) > 0 {
		var replyIdList []string
		for id := range replyIds {
			replyIdList = append(replyIdList, id)
		}
		var replyTargets []moment_models.MomentCommentModel
		l.svcCtx.DB.Where("comment_id IN (?)", replyIdList).Find(&replyTargets)
		for _, rt := range replyTargets {
			replyTargetMap[rt.CommentID] = rt
			userIds[rt.UserID] = true
		}
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

	// 转换评论数据（顶层 + 预览 children）
	var commentInfos []types.GetMomentCommentsInfo
	for _, comment := range topComments {
		userInfo := userInfoMap[comment.UserID]
		userName := ""
		avatar := ""
		if userInfo != nil {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		var childInfos []types.GetMomentCommentsInfo
		for _, ch := range childMap[comment.CommentID] {
			chUser := userInfoMap[ch.UserID]
			chName := ""
			chAvatar := ""
			if chUser != nil {
				chName = chUser.NickName
				chAvatar = chUser.Avatar
			}
			replyName := ""
			if ch.ReplyToCommentID != "" {
				if target, ok := replyTargetMap[ch.ReplyToCommentID]; ok {
					if u := userInfoMap[target.UserID]; u != nil {
						replyName = u.NickName
					}
				}
			}
			childInfos = append(childInfos, types.GetMomentCommentsInfo{
				Id:               ch.CommentID,
				UserID:           ch.UserID,
				UserName:         chName,
				Avatar:           chAvatar,
				Content:          ch.Content,
				ParentId:         ch.ParentID,
				ReplyToCommentId: ch.ReplyToCommentID,
				ReplyToUserName:  replyName,
				ChildCount:       0,
				Children:         nil,
				CreatedAt:        ch.CreatedAt.String(),
			})
		}
		if len(childInfos) > childPreviewLimit {
			childInfos = childInfos[:childPreviewLimit]
		}

		commentInfos = append(commentInfos, types.GetMomentCommentsInfo{
			Id:               comment.CommentID,
			UserID:           comment.UserID,
			UserName:         userName,
			Avatar:           avatar,
			Content:          comment.Content,
			ParentId:         comment.ParentID,
			ReplyToCommentId: comment.ReplyToCommentID,
			ReplyToUserName:  "",
			ChildCount:       childCountMap[comment.CommentID],
			Children:         childInfos,
			CreatedAt:        comment.CreatedAt.String(),
		})
	}

	// 构建响应
	resp = &types.GetMomentCommentsRes{
		Count: totalCount,
		List:  commentInfos,
	}

	return resp, nil
}
