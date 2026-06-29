package post

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePostLogic) CreatePost(req *types.CreatePostReq) (resp *types.CreatePostRes, err error) {
	// 校验是否是圈子成员
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("请先加入圈子再发帖")
	}

	postID := uuid.New().String()
	var files *circle_models.PostFiles
	if len(req.Files) > 0 {
		f := make(circle_models.PostFiles, 0, len(req.Files))
		for _, fi := range req.Files {
			f = append(f, circle_models.PostFileInfo{FileKey: fi.FileKey, Type: fi.Type})
		}
		files = &f
	}

	post := circle_models.CirclePostModel{
		PostID:   postID,
		CircleID: req.CircleID,
		UserID:   req.UserID,
		Title:    req.Title,
		Content:  req.Content,
		Files:    files,
	}
	if err = l.svcCtx.DB.Create(&post).Error; err != nil {
		return nil, fmt.Errorf("发布帖子失败: %v", err)
	}

	// 更新圈子帖子数
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", req.CircleID).
		UpdateColumn("post_count", l.svcCtx.DB.Raw("post_count + 1"))

	userName, avatar := "", ""
	userResp, _ := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: []string{req.UserID}})
	if userResp != nil {
		if info := userResp.UserInfo[req.UserID]; info != nil {
			userName = info.NickName
			avatar = info.Avatar
		}
	}

	return &types.CreatePostRes{
		PostID:    postID,
		CircleID:  req.CircleID,
		UserID:    req.UserID,
		UserName:  userName,
		Avatar:    avatar,
		Title:     req.Title,
		Content:   req.Content,
		CreatedAt: post.CreatedAt.String(),
	}, nil
}
