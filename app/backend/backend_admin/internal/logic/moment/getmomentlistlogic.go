package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/moment/moment_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMomentListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取动态列表
func NewGetMomentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMomentListLogic {
	return &GetMomentListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMomentListLogic) GetMomentList(req *types.GetMomentListReq) (resp *types.GetMomentListRes, err error) {
	// 分页参数校验
	page := req.Page
	limit := req.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 用户ID筛选
	if req.UserId != "" {
		whereClause = whereClause.Where("user_id = ?", req.UserId)
	}

	// 分页查询
	moments, count, err := list_query.ListQuery(l.svcCtx.DB, moment_models.MomentModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  page,
			Limit: limit,
			Key:   req.Keywords,
			Sort:  "created_at desc",
		},
		Where: whereClause,
		Likes: []string{"content"},
	})

	if err != nil {
		logx.Errorf("查询动态列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.GetMomentListItem
	for _, moment := range moments {
		var files []types.GetMomentListFileInfo
		if moment.Files != nil {
			for _, file := range *moment.Files {
				files = append(files, types.GetMomentListFileInfo{
					FileName: file.FileKey,
				})
			}
		}

		list = append(list, types.GetMomentListItem{
			Id:        moment.UUID,
			UserId:    moment.UserID,
			Content:   moment.Content,
			Files:     files,
			IsDeleted: moment.IsDeleted,
			CreatedAt: moment.CreatedAt.String(),
			UpdatedAt: moment.UpdatedAt.String(),
		})
	}

	return &types.GetMomentListRes{
		List:  list,
		Total: count,
	}, nil
}
