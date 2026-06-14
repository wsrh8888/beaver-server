package logic

import (
	"context"
	"errors"
	"time"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ListUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUsersLogic {
	return &ListUsersLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListUsersLogic) ListUsers(in *user_rpc.ListUsersReq) (*user_rpc.ListUsersRes, error) {
	if in.UserId != "" {
		var user user_models.UserModel
		if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).First(&user).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &user_rpc.ListUsersRes{}, nil
			}
			return nil, err
		}
		return &user_rpc.ListUsersRes{
			Total: 1,
			List:  []*user_rpc.UserInfo{userInfoFromModel(user)},
		}, nil
	}

	page := int(in.Page)
	if page <= 0 {
		page = 1
	}
	pageSize := int(in.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	db := l.svcCtx.DB.Model(&user_models.UserModel{})
	if in.Status > 0 {
		db = db.Where("status = ?", in.Status)
	}
	if in.Source > 0 {
		db = db.Where("source = ?", in.Source)
	}
	if in.UserType > 0 {
		db = db.Where("user_type = ?", in.UserType)
	}
	if in.Email != "" {
		db = db.Where("email LIKE ?", "%"+in.Email+"%")
	}
	if in.Keyword != "" {
		kw := "%" + in.Keyword + "%"
		db = db.Where("nick_name LIKE ? OR email LIKE ? OR phone LIKE ?", kw, kw, kw)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		l.Errorf("统计用户失败: %v", err)
		return nil, err
	}

	var users []user_models.UserModel
	if err := db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		l.Errorf("查询用户列表失败: %v", err)
		return nil, err
	}

	list := make([]*user_rpc.UserInfo, 0, len(users))
	for _, user := range users {
		list = append(list, userInfoFromModel(user))
	}
	return &user_rpc.ListUsersRes{Total: total, List: list}, nil
}

func userInfoFromModel(user user_models.UserModel) *user_rpc.UserInfo {
	return &user_rpc.UserInfo{
		UserId:    user.UserID,
		NickName:  user.NickName,
		Avatar:    user.Avatar,
		Version:   user.Version,
		Email:     user.Email,
		Abstract:  user.Abstract,
		Phone:     user.Phone,
		Status:    int32(user.Status),
		Source:    user.Source,
		UserType:  int32(user.UserType),
		CreatedAt: time.Time(user.CreatedAt).Format(time.RFC3339),
		UpdatedAt: time.Time(user.UpdatedAt).Format(time.RFC3339),
	}
}
