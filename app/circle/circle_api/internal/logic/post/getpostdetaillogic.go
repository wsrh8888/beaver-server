package post

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostDetailLogic {
	return &GetPostDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostDetailLogic) GetPostDetail(req *types.GetPostDetailReq) (resp *types.GetPostDetailRes, err error) {
	var p circle_models.CirclePostModel
	if err = l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", req.PostID).First(&p).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	// 用户信息
	userName, avatar := "", ""
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{p.UserID}})
	if userResp != nil {
		if info := userResp.UserInfo[p.UserID]; info != nil {
			userName = info.NickName
			avatar = info.Avatar
		}
	}

	// 是否点赞
	isLiked := false
	var like circle_models.CircleLikeModel
	if l.svcCtx.DB.Where("post_id = ? AND user_id = ?", req.PostID, req.UserID).First(&like).Error == nil {
		isLiked = true
	}

	// 最新20条一级评论
	var comments []circle_models.CircleCommentModel
	l.svcCtx.DB.Where("post_id = ? AND parent_id = '' AND is_deleted = false", req.PostID).
		Order("created_at DESC").Limit(20).Find(&comments)

	commentItems := buildCommentItems(l.ctx, l.svcCtx, comments, req.PostID)

	resp = &types.GetPostDetailRes{
		PostID:       p.PostID,
		CircleID:     p.CircleID,
		UserID:       p.UserID,
		UserName:     userName,
		Avatar:       avatar,
		Title:        p.Title,
		Content:      p.Content,
		CommentCount: p.CommentCount,
		LikeCount:    p.LikeCount,
		IsLiked:      isLiked,
		IsTop:        p.IsTop,
		Comments:     commentItems,
		CreatedAt:    p.CreatedAt.String(),
	}
	if p.Files != nil {
		for _, f := range *p.Files {
			resp.Files = append(resp.Files, types.GetPostDetailFileInfo{FileKey: f.FileKey, Type: f.Type})
		}
	}
	return resp, nil
}

func buildCommentItems(ctx context.Context, svcCtx *svc.ServiceContext, comments []circle_models.CircleCommentModel, postID string) []types.GetPostDetailCommentInfo {
	if len(comments) == 0 {
		return []types.GetPostDetailCommentInfo{}
	}

	// 收集所有评论用户ID
	userIDs := make([]string, 0)
	commentIDs := make([]string, 0)
	for _, c := range comments {
		userIDs = append(userIDs, c.UserID)
		commentIDs = append(commentIDs, c.CommentID)
	}

	userResp, _ := svcCtx.UserRpc.UserListInfo(ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	// 查子评论数
	type countResult struct {
		ParentID string
		Count    int64
	}
	var counts []countResult
	svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Select("parent_id, count(*) as count").
		Where("parent_id IN ? AND is_deleted = false", commentIDs).
		Group("parent_id").
		Scan(&counts)
	childCountMap := make(map[string]int64)
	for _, cr := range counts {
		childCountMap[cr.ParentID] = cr.Count
	}

	items := make([]types.GetPostDetailCommentInfo, 0, len(comments))
	for _, c := range comments {
		item := types.GetPostDetailCommentInfo{
			CommentID:        c.CommentID,
			UserID:           c.UserID,
			Content:          c.Content,
			ParentID:         c.ParentID,
			ReplyToCommentID: c.ReplyToCommentID,
			ChildCount:       childCountMap[c.CommentID],
			Children:         []types.GetPostDetailCommentInfo{},
			CreatedAt:        c.CreatedAt.String(),
		}
		if userResp != nil {
			if info := userResp.UserInfo[c.UserID]; info != nil {
				item.UserName = info.NickName
				item.Avatar = info.Avatar
			}
			if c.ReplyToUserID != "" {
				if info := userResp.UserInfo[c.ReplyToUserID]; info != nil {
					item.ReplyToUserName = info.NickName
				}
			}
		}
		items = append(items, item)
	}
	return items
}
