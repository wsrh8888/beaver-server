package logic

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePostLogic) DeletePost(in *circle_rpc.DeletePostReq) (*circle_rpc.DeletePostRes, error) {
	var post circle_models.CirclePostModel
	if err := l.svcCtx.DB.Where("post_id = ? AND is_deleted = false", in.PostId).First(&post).Error; err != nil {
		return nil, fmt.Errorf("帖子不存在")
	}

	if err := l.svcCtx.DB.Model(&post).Update("is_deleted", true).Error; err != nil {
		return nil, fmt.Errorf("删除帖子失败: %v", err)
	}

	return &circle_rpc.DeletePostRes{}, nil
}
