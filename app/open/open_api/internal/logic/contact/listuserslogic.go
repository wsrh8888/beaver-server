package contact

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 分页查询用户列表
func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListUsersLogic) ListUsers(req *types.ListUsersReq) (resp *types.ListUsersRes, err error) {
	// 1. 构建查询
	query := l.svcCtx.DB.Model(&user_models.UserModel{})

	// 2. 添加过滤条件
	if req.Keyword != "" {
		query = query.Where("nick_name LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	if req.Department != "" {
		// TODO: 需要根据实际的部门关联表来查询
		// query = query.Where("department = ?", req.Department)
	}
	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 3. 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// 4. 分页查询
	var users []user_models.UserModel
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&users).Error; err != nil {
		return nil, err
	}

	// 5. 转换为响应格式
	var userList []types.UserDetailInfo
	for _, user := range users {
		userList = append(userList, types.UserDetailInfo{
			UserID:   user.UserID,
			Nickname: user.NickName,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
			Email:    user.Email,
		})
	}

	return &types.ListUsersRes{
		Total: total,
		Users: userList,
	}, nil
}
