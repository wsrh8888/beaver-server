package oauth

import (
	"context"
	"errors"

	"beaver/app/open/open_api/internal/logic/oauthutil"
	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetH5AuthCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetH5AuthCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetH5AuthCodeLogic {
	return &GetH5AuthCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetH5AuthCodeLogic) GetH5AuthCode(req *types.GetH5AuthCodeReq) (resp *types.GetH5AuthCodeRes, err error) {
	if req.UserID == "" {
		return nil, errors.New("未登录")
	}
	if req.AppID == "" {
		return nil, errors.New("appId 不能为空")
	}

	authCode, expireIn, err := oauthutil.CreateOAuthCode(l.svcCtx.DB, req.AppID, req.UserID, "h5_sso")
	if err != nil {
		logx.Errorf("生成 H5 authCode 失败: appId=%s, userId=%s, err=%v", req.AppID, req.UserID, err)
		return nil, err
	}

	logx.Infof("生成 H5 authCode 成功: appId=%s, userId=%s", req.AppID, req.UserID)

	return &types.GetH5AuthCodeRes{
		AuthCode: authCode,
		ExpireIn: expireIn,
	}, nil
}
