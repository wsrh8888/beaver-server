package comment

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCommentLogic {
	return &DeleteCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteCommentLogic) DeleteComment(req *types.DeleteCommentReq) (resp *types.DeleteCommentRes, err error) {
	var c circle_models.CircleCommentModel
	if err = l.svcCtx.DB.Where("comment_id = ? AND is_deleted = false", req.CommentID).First(&c).Error; err != nil {
		return nil, fmt.Errorf("评论不存在")
	}

	// 评论人可删，圈主/管理员也可删
	if c.UserID != req.UserID {
		var member circle_models.CircleMemberModel
		if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", c.CircleID, req.UserID).First(&member).Error; err != nil {
			return nil, fmt.Errorf("无权限删除")
		}
		if member.Role > 2 {
			return nil, fmt.Errorf("无权限删除")
		}
	}

	if err = l.svcCtx.DB.Model(&c).Update("is_deleted", true).Error; err != nil {
		return nil, fmt.Errorf("删除评论失败: %v", err)
	}

	return &types.DeleteCommentRes{}, nil
}
