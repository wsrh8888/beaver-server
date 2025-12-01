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

	// 获取评论总数
	var totalCount int64
	if err := l.svcCtx.DB.Model(&moment_models.MomentCommentModel{}).
		Where("moment_id = ? AND is_deleted = false", req.MomentID).
		Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// 获取分页评论数据
	var comments []moment_models.MomentCommentModel
	if err := l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false", req.MomentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	// 获取评论用户ID列表
	userIds := make(map[string]bool)
	for _, comment := range comments {
		userIds[comment.UserID] = true
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

	// 转换评论数据
	var commentInfos []types.GetMomentCommentsInfo
	for _, comment := range comments {
		userInfo := userInfoMap[comment.UserID]
		userName := ""
		avatar := ""
		if userInfo != nil {
			userName = userInfo.NickName
			avatar = userInfo.Avatar
		}

		commentInfos = append(commentInfos, types.GetMomentCommentsInfo{
			Id:        comment.UUID,
			UserID:    comment.UserID,
			UserName:  userName,
			Avatar:    avatar,
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt.String(),
		})
	}

	// 构建响应
	resp = &types.GetMomentCommentsRes{
		Count: totalCount,
		List:  commentInfos,
	}

	return resp, nil
}
