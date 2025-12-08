package logic

import (
	"context"
	"errors"

	"beaver/app/user/user_models"
	"beaver/app/user/user_rpc/internal/svc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserLogic) SearchUser(in *user_rpc.SearchUserReq) (*user_rpc.SearchUserRes, error) {
	var user user_models.UserModel

	var err error
	switch in.Type {
	case "email":
		err = l.svcCtx.DB.Take(&user, "email = ?", in.Keyword).Error
	case "userId":
		err = l.svcCtx.DB.Take(&user, "user_id = ?", in.Keyword).Error
	default:
		// 默认按邮箱搜索
		err = l.svcCtx.DB.Take(&user, "email = ?", in.Keyword).Error
	}

	if err != nil {
		l.Logger.Errorf("搜索用户失败: keyword=%s, type=%s, error=%v", in.Keyword, in.Type, err)
		return nil, errors.New("用户不存在")
	}

	return &user_rpc.SearchUserRes{
		UserInfo: &user_rpc.UserInfo{
			UserId:   user.UserID,
			NickName: user.NickName,
			Avatar:   user.Avatar,
			Version:  user.Version,
			Email:    user.Email,
		},
	}, nil
}
