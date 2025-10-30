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
		Preload("CommentsModel.CommentUserModel").
		Preload("LikesModel.LikeUserModel").
		Preload("MomentUserModel").
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
		Count: count, // 修改为动态数目
		List:  make([]types.MomentModel, 0, len(moments)),
	}

	for _, moment := range moments {
		var likes []types.MomentLikeModel
		for _, like := range moment.LikesModel {
			likes = append(likes, types.MomentLikeModel{
				Id:       like.Id,
				UserID:   like.UserID,
				UserName: like.LikeUserModel.NickName,
				Avatar:   like.LikeUserModel.Avatar,
			})
		}

		var files []types.FileInfo
		for _, file := range *moment.Files {
			files = append(files, types.FileInfo{
				FileKey: file.FileKey,
			})
		}

		resp.List = append(resp.List, types.MomentModel{
			Id:      moment.Id,
			UserID:  moment.MomentUserModel.UUID, // 使用 MomentUserModel 的 UserID
			Content: moment.Content,
			Files:   files,
			Likes:   likes,
			// Comments: comments,
			UserName:  moment.MomentUserModel.NickName, // 增加 UserName 字段
			Avatar:    moment.MomentUserModel.Avatar,   // 增加 Avatar 字段
			CreatedAt: moment.CreatedAt.String(),
		})
	}

	return resp, nil
}
