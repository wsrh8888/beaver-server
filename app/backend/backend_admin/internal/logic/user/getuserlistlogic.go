package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/user/user_models"
	"beaver/common/list_query"
	"beaver/common/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取用户列表
func NewGetUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserListLogic {
	return &GetUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserListLogic) GetUserList(req *types.GetUserListReq) (resp *types.GetUserListRes, err error) {
	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 构建查询条件
	whereClause := l.svcCtx.DB.Where("1 = 1")

	// 状态筛选
	if req.Status != 0 {
		whereClause = whereClause.Where("status = ?", req.Status)
	}

	// 来源筛选
	if req.Source != 0 {
		whereClause = whereClause.Where("source = ?", req.Source)
	}

	// 邮箱筛选
	if req.Email != "" {
		whereClause = whereClause.Where("email LIKE ?", "%"+req.Email+"%")
	}

	// 分页查询
	users, count, err := list_query.ListQuery(l.svcCtx.DB, user_models.UserModel{}, list_query.Option{
		PageInfo: models.PageInfo{
			Page:  req.Page,
			Limit: req.PageSize,
			Sort:  "created_at desc",
		},
		Where: whereClause,
	})

	if err != nil {
		l.Logger.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	// 转换为响应格式
	var list []types.UserInfo
	for _, user := range users {
		list = append(list, types.UserInfo{
			Id:          user.UUID,
			Nickname:    user.NickName,
			Email:       user.Email,
			Abstract:    user.Abstract,
			FileName:    user.Avatar,
			Status:      int(user.Status),
			Source:      int(user.Source),
			LastLoginIP: "", // UserModel 没有 LastLoginIP 字段
			CreateTime:  user.CreatedAt.String(),
			UpdateTime:  user.UpdatedAt.String(),
		})
	}

	l.Logger.Infof("获取用户列表成功: page=%d, pageSize=%d, total=%d", req.Page, req.PageSize, count)
	return &types.GetUserListRes{
		List:  list,
		Total: count,
	}, nil
}
