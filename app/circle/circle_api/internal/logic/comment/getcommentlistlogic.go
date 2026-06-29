package comment

import (
	"context"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCommentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentListLogic {
	return &GetCommentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCommentListLogic) GetCommentList(req *types.GetCommentListReq) (resp *types.GetCommentListRes, err error) {
	var total int64
	var comments []circle_models.CircleCommentModel

	query := l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Where("post_id = ? AND is_deleted = false", req.PostID)
	if req.ParentID != "" {
		query = query.Where("parent_id = ?", req.ParentID)
	} else {
		query = query.Where("parent_id = ''")
	}
	query.Count(&total)
	query.Order("created_at ASC").
		Offset((req.Page - 1) * req.Limit).
		Limit(req.Limit).
		Find(&comments)

	if len(comments) == 0 {
		return &types.GetCommentListRes{Count: total, List: []types.GetCommentListItem{}}, nil
	}

	// 批量拉用户信息
	userIDs := make([]string, 0, len(comments))
	commentIDs := make([]string, 0, len(comments))
	for _, c := range comments {
		userIDs = append(userIDs, c.UserID)
		if c.ReplyToUserID != "" {
			userIDs = append(userIDs, c.ReplyToUserID)
		}
		commentIDs = append(commentIDs, c.CommentID)
	}
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs})

	// 查子评论数
	type countResult struct {
		ParentID string
		Count    int64
	}
	var counts []countResult
	l.svcCtx.DB.Model(&circle_models.CircleCommentModel{}).
		Select("parent_id, count(*) as count").
		Where("parent_id IN ? AND is_deleted = false", commentIDs).
		Group("parent_id").
		Scan(&counts)
	childCountMap := make(map[string]int64)
	for _, cr := range counts {
		childCountMap[cr.ParentID] = cr.Count
	}

	items := make([]types.GetCommentListItem, 0, len(comments))
	for _, c := range comments {
		item := types.GetCommentListItem{
			CommentID:        c.CommentID,
			UserID:           c.UserID,
			Content:          c.Content,
			ParentID:         c.ParentID,
			ReplyToCommentID: c.ReplyToCommentID,
			ChildCount:       childCountMap[c.CommentID],
			Children:         []types.GetCommentListItem{},
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

	return &types.GetCommentListRes{Count: total, List: items}, nil
}
