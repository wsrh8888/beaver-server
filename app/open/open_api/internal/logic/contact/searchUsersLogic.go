package contact

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"
	user_models "beaver/app/user/user_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUsersLogic {
	return &SearchUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchUsersLogic) SearchUsers(req *types.SearchUsersReq) (resp *types.SearchUsersRes, err error) {
	// 1. 参数校验
	if req.Keyword == "" {
		return nil, errors.New("搜索关键词不能为空")
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// 2. 查询总数
	var total int64
	err = l.svcCtx.DB.Model(&user_models.UserModel{}).
		Where("nick_name LIKE ? OR phone LIKE ? OR email LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%").
		Count(&total).Error
	if err != nil {
		return nil, errors.New("查询失败")
	}

	// 3. 分页查询
	var users []user_models.UserModel
	offset := (req.Page - 1) * req.PageSize
	err = l.svcCtx.DB.Where("nick_name LIKE ? OR phone LIKE ? OR email LIKE ?",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%").
		Offset(offset).Limit(req.PageSize).
		Find(&users).Error
	if err != nil {
		return nil, errors.New("查询失败")
	}

	// 4. 转换为响应格式
	var userDetails []types.UserDetail
	for _, user := range users {
		userDetails = append(userDetails, types.UserDetail{
			UserID:   user.UserID,
			NickName: user.NickName,
			Avatar:   user.Avatar,
			Phone:    user.Phone,
			Email:    user.Email,
			Gender:   int(user.Gender),
			Status:   int(user.Status),
		})
	}

	return &types.SearchUsersRes{
		Total: total,
		Users: userDetails,
	}, nil
}
