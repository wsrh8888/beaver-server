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
	// 定义 MomentModel 实例
	var moment moment_models.MomentModel

	// 查询数据库并预加载关联信息
	if err := l.svcCtx.DB.
		Preload("Comments").
		Preload("Likes").
		First(&moment, "id = ?", req.MomentID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("moment not found")
		}
		return nil, err
	}

	// 构造响应
	resp = &types.GetMomentInfoRes{
		Moment: convertMomentModel(moment),
	}

	return resp, nil
}

func convertMomentModel(moment moment_models.MomentModel) types.MomentModel {
	return types.MomentModel{
		Id:      moment.Id,
		UserID:  moment.UserID,
		Content: moment.Content,
		Files:   convertToResponseFiles(*moment.Files),
		Likes:   convertToResponseLikes(moment.LikesModel),
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
