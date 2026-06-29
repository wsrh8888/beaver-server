package post

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostListLogic {
	return &GetPostListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPostListLogic) GetPostList(req *types.GetPostListReq) (resp *types.GetPostListRes, err error) {
	var total int64
	var posts []circle_models.CirclePostModel

	l.svcCtx.DB.Model(&circle_models.CirclePostModel{}).
		Where("circle_id = ? AND is_deleted = false", req.CircleID).
		Count(&total)
	l.svcCtx.DB.Where("circle_id = ? AND is_deleted = false", req.CircleID).
		Order("is_top DESC, created_at DESC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&posts)

	if len(posts) == 0 {
		return &types.GetPostListRes{Count: total, List: []types.PostListItem{}}, nil
	}

	// 批量拉用户信息
	userIDs := make([]string, 0, len(posts))
	postIDs := make([]string, 0, len(posts))
	for _, p := range posts {
		userIDs = append(userIDs, p.UserID)
		postIDs = append(postIDs, p.PostID)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	// 查当前用户点赞情况
	var likes []circle_models.CircleLikeModel
	l.svcCtx.DB.Where("post_id IN ? AND user_id = ?", postIDs, req.UserID).Find(&likes)
	likedMap := make(map[string]bool)
	for _, lk := range likes {
		likedMap[lk.PostID] = true
	}

	// 每帖取最新3条评论
	commentMap := make(map[string][]circle_models.CircleCommentModel)
	var allComments []circle_models.CircleCommentModel
	l.svcCtx.DB.Where("post_id IN ? AND parent_id = '' AND is_deleted = false", postIDs).
		Order("created_at DESC").
		Find(&allComments)
	countMap := make(map[string]int)
	for _, c := range allComments {
		if countMap[c.PostID] < 3 {
			commentMap[c.PostID] = append(commentMap[c.PostID], c)
			countMap[c.PostID]++
		}
	}

	items := make([]types.PostListItem, 0, len(posts))
	for _, p := range posts {
		item := types.PostListItem{
			PostID:       p.PostID,
			CircleID:     p.CircleID,
			UserID:       p.UserID,
			Title:        p.Title,
			Content:      p.Content,
			CommentCount: p.CommentCount,
			LikeCount:    p.LikeCount,
			IsLiked:      likedMap[p.PostID],
			IsTop:        p.IsTop,
			CreatedAt:    p.CreatedAt.String(),
		}
		if userResp != nil {
			if info := userResp.UserInfo[p.UserID]; info != nil {
				item.UserName = info.NickName
				item.Avatar = info.Avatar
			}
		}
		if p.Files != nil {
			for _, f := range *p.Files {
				item.Files = append(item.Files, types.GetPostListFileInfo{FileKey: f.FileKey, Type: f.Type})
			}
		}
		for _, c := range commentMap[p.PostID] {
			item.Comments = append(item.Comments, types.GetPostListCommentInfo{
				CommentID: c.CommentID,
				UserID:    c.UserID,
				Content:   c.Content,
				CreatedAt: c.CreatedAt.String(),
			})
		}
		items = append(items, item)
	}

	return &types.GetPostListRes{Count: total, List: items}, nil
}
