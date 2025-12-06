package logic

import (
	"context"
	"fmt"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetMomentDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态详情的接口（包含更多评论和点赞）
func NewGetMomentDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentDetailLogic {
	return &GetMomentDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentDetailLogic) GetMomentDetail(req *types.GetMomentDetailReq) (resp *types.GetMomentDetailRes, err error) {
	// 获取动态信息
	var moment moment_models.MomentModel
	if err := l.svcCtx.DB.Where("uuid = ? AND is_deleted = false", req.MomentID).First(&moment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("moment not found")
		}
		return nil, err
	}

	// 获取动态用户信息
	userResp, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{
		UserIdList: []string{moment.UserID},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}

	userInfo := userResp.UserInfo[moment.UserID]
	if userInfo == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 获取评论数量
	var commentCount int64
	if err := l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
		Where("moment_id = ? AND is_deleted = false", req.MomentID).
		Count(&commentCount).Error; err != nil {
		return nil, err
	}

	// 获取点赞数量
	var likeCount int64
	if err := l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
		Where("moment_id = ? AND is_deleted = false", req.MomentID).
		Count(&likeCount).Error; err != nil {
		return nil, err
	}

	// 检查当前用户是否已点赞
	var isLiked bool
	if req.UserID != "" {
		var likeExists int64
		if err := l.svcCtx.DB.Model(&moment_models.MomentLikeModel{}).
			Where("moment_id = ? AND user_id = ? AND is_deleted = false", req.MomentID, req.UserID).
			Limit(1).
			Count(&likeExists).Error; err == nil && likeExists > 0 {
			isLiked = true
		}
	}

	// 获取最新的20条顶层评论
	var topComments []moment_models.MomentCommentModel
	if commentCount > 0 {
		if err := l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false AND parent_id = ''", req.MomentID).
			Order("created_at DESC").
			Limit(20).
			Find(&topComments).Error; err != nil {
			return nil, err
		}
	}

	// 统计子评论数量
	childCountMap := make(map[string]int64)
	if len(topComments) > 0 {
		var childStats []struct {
			ParentID string
			Count    int64
		}
		var parentIDs []string
		for _, c := range topComments {
			parentIDs = append(parentIDs, c.UUID)
		}
		l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
			Where("parent_id IN (?) AND is_deleted = false", parentIDs).
			Select("parent_id, COUNT(*) as count").
			Group("parent_id").
			Scan(&childStats)
		for _, stat := range childStats {
			childCountMap[stat.ParentID] = stat.Count
		}
	}

	// 获取子评论预览（每个顶层最多20条）
	childPreviewLimit := 20
	childMap := make(map[string][]moment_models.MomentCommentModel)
	if len(topComments) > 0 {
		var parentIDs []string
		for _, c := range topComments {
			parentIDs = append(parentIDs, c.UUID)
		}
		var children []moment_models.MomentCommentModel
		l.svcCtx.DB.Where("parent_id IN (?) AND is_deleted = false", parentIDs).
			Order("created_at DESC").
			Find(&children)
		childCount := make(map[string]int)
		for _, ch := range children {
			if childCount[ch.ParentID] < childPreviewLimit {
				childMap[ch.ParentID] = append(childMap[ch.ParentID], ch)
				childCount[ch.ParentID]++
			}
		}
	}

	// 获取最新的50个点赞
	var likes []moment_models.MomentLikeModel
	if likeCount > 0 {
		if err := l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false", req.MomentID).
			Order("created_at DESC").
			Limit(50).
			Find(&likes).Error; err != nil {
			return nil, err
		}
	}

	// 获取评论和点赞的用户信息
	userIds := make(map[string]bool)
	userIds[moment.UserID] = true // 动态作者
	for _, comment := range topComments {
		userIds[comment.UserID] = true
	}
	for _, children := range childMap {
		for _, c := range children {
			userIds[c.UserID] = true
		}
	}
	for _, like := range likes {
		userIds[like.UserID] = true
	}

	var userIdList []string
	for userId := range userIds {
		userIdList = append(userIdList, userId)
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

	// 转换文件信息
	var files []types.GetMomentDetailFileInfo
	if moment.Files != nil {
		for _, file := range *moment.Files {
			files = append(files, types.GetMomentDetailFileInfo{
				FileKey: file.FileKey,
				Type:    uint32(file.Type),
			})
		}
	}

	// 转换评论信息
	var commentInfos []types.GetMomentDetailCommentInfo
	for _, comment := range topComments {
		userInfo := userInfoMap[comment.UserID]
		userName := ""
		avatar := ""
		if userInfo != nil {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		// 子评论预览
		var childInfos []types.GetMomentDetailCommentInfo
		for _, ch := range childMap[comment.UUID] {
			chUser := userInfoMap[ch.UserID]
			chName := ""
			chAvatar := ""
			if chUser != nil {
				chName = chUser.NickName
				chAvatar = chUser.Avatar
			}

			// 需要被回复的昵称
			replyName := ""
			if ch.ReplyToCommentID != "" {
				// 尝试从子列表或父评论获取
				if ch.ReplyToCommentID == comment.UUID {
					replyName = userName
				} else {
					// 在当前 childMap 里查找
					for _, other := range childMap[comment.UUID] {
						if other.UUID == ch.ReplyToCommentID {
							if u := userInfoMap[other.UserID]; u != nil {
								replyName = u.NickName
							}
							break
						}
					}
				}
			}

			childInfos = append(childInfos, types.GetMomentDetailCommentInfo{
				Id:               ch.UUID,
				UserID:           ch.UserID,
				UserName:         chName,
				Avatar:           chAvatar,
				Content:          ch.Content,
				ParentId:         ch.ParentID,
				ReplyToCommentId: ch.ReplyToCommentID,
				ReplyToUserName:  replyName,
				ChildCount:       0, // 子评论不再嵌套
				Children:         nil,
				CreatedAt:        ch.CreatedAt.String(),
			})
		}

		commentInfos = append(commentInfos, types.GetMomentDetailCommentInfo{
			Id:               comment.UUID,
			UserID:           comment.UserID,
			UserName:         userName,
			Avatar:           avatar,
			Content:          comment.Content,
			ParentId:         "", // 顶层
			ReplyToCommentId: "",
			ReplyToUserName:  "",
			ChildCount:       childCountMap[comment.UUID],
			Children:         childInfos,
			CreatedAt:        comment.CreatedAt.String(),
		})
	}

	// 转换点赞信息
	var likeInfos []types.GetMomentDetailLikeInfo
	for _, like := range likes {
		userInfo := userInfoMap[like.UserID]
		userName := ""
		avatar := ""
		if userInfo != nil {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		likeInfos = append(likeInfos, types.GetMomentDetailLikeInfo{
			Id:        like.UUID,
			UserID:    like.UserID,
			UserName:  userName,
			Avatar:    avatar,
			CreatedAt: like.CreatedAt.String(),
		})
	}

	// 构建响应
	resp = &types.GetMomentDetailRes{
		Id:           moment.UUID,
		UserID:       moment.UserID,
		UserName:     userInfo.NickName,
		Avatar:       userInfo.Avatar,
		Content:      moment.Content,
		Files:        files,
		Comments:     commentInfos,
		Likes:        likeInfos,
		CommentCount: commentCount,
		LikeCount:    likeCount,
		IsLiked:      isLiked,
		CreatedAt:    moment.CreatedAt.String(),
	}

	return resp, nil
}
