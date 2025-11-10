package logic

import (
	"context"
	"fmt"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetMomentInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMomentInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentInfoLogic {
	return &GetMomentInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentInfoLogic) GetMomentInfo(req *types.GetMomentInfoReq) (resp *types.GetMomentInfoRes, err error) {
	// 查询动态基本信息
	var moment moment_models.MomentModel
	if err := l.svcCtx.DB.First(&moment, "id = ?", req.MomentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("moment not found")
		}
		return nil, err
	}

	// 查询评论信息
	var comments []moment_models.MomentCommentModel
	l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false", req.MomentID).Find(&comments)

	// 查询点赞信息
	var likes []moment_models.MomentLikeModel
	l.svcCtx.DB.Where("moment_id = ? AND is_deleted = false", req.MomentID).Find(&likes)

	// 构造响应
	resp = &types.GetMomentInfoRes{
		Moment: convertMomentModel(moment, comments, likes),
	}

	return resp, nil
}

func convertMomentModel(moment moment_models.MomentModel, comments []moment_models.MomentCommentModel, likes []moment_models.MomentLikeModel) types.MomentModel {
	return types.MomentModel{
		Id:       moment.Id,
		UserID:   moment.UserID,
		Content:  moment.Content,
		Files:    convertToResponseFiles(*moment.Files),
		Comments: convertToResponseComments(comments),
		Likes:    convertToResponseLikes(likes),
	}
}

func convertToResponseFiles(files moment_models.Files) []types.FileInfo {
	var result []types.FileInfo
	for _, file := range files {
		result = append(result, types.FileInfo{
			FileKey: file.FileKey,
		})
	}
	return result
}

func convertToResponseComments(comments []moment_models.MomentCommentModel) []types.MomentCommentModel {
	var result []types.MomentCommentModel
	for _, comment := range comments {
		result = append(result, types.MomentCommentModel{
			Id:       comment.Id,
			MomentId: comment.MomentID,
			UserID:   comment.UserID,
			Content:  comment.Content,
		})
	}
	return result
}

func convertToResponseLikes(likes []moment_models.MomentLikeModel) []types.MomentLikeModel {
	var result []types.MomentLikeModel
	for _, like := range likes {
		result = append(result, types.MomentLikeModel{
			Id:       like.Id,
			MomentId: like.MomentID,
			UserID:   like.UserID,
		})
	}
	return result
}
