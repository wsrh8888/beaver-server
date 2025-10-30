package logic

import (
	"context"
	"fmt"

	"beaver/app/moment/moment_api/internal/svc"
	"beaver/app/moment/moment_api/internal/types"
	"beaver/app/moment/moment_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMomentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateMomentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateMomentLogic {
	return &CreateMomentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateMomentLogic) CreateMoment(req *types.CreateMomentReq) (resp *types.CreateMomentRes, err error) {
	// 构造MomentModel实例
	moment := moment_models.MomentModel{
		UserID:        req.UserID,
		Content:       req.Content,
		Files:         convertFiles(req.Files),
		CommentsModel: []moment_models.MomentCommentModel{}, // 初始化为空切片
		LikesModel:    []moment_models.MomentLikeModel{},    // 初始化为空切片
	}

	// 插入数据库
	if err := l.svcCtx.DB.Create(&moment).Error; err != nil {
		return nil, fmt.Errorf("failed to create moment: %v", err)
	}

	return resp, nil
}

// 辅助函数：将请求中的文件信息转换为数据库模型所需的结构
func convertFiles(files []types.FileInfo) *moment_models.Files {
	var result moment_models.Files
	for _, file := range files {
		result = append(result, moment_models.FileInfo{
			FileKey: file.FileKey,
		})
	}
	return &result
}
