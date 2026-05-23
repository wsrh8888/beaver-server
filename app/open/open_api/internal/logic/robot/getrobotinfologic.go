package robot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRobotInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Robot 自身信息（AppID+AppSecret 换到 token 后调此接口确认 Robot 身份）
func NewGetRobotInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRobotInfoLogic {
	return &GetRobotInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRobotInfoLogic) GetRobotInfo(req *types.GetRobotInfoReq) (resp *types.GetRobotInfoRes, err error) {
	// todo: add your logic here and delete this line

	return
}
