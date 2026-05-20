package bot

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBotInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取 Bot 自身信息（AppID+AppSecret 换到 token 后调此接口确认 Bot 身份）
func NewGetBotInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBotInfoLogic {
	return &GetBotInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBotInfoLogic) GetBotInfo(req *types.GetBotInfoReq) (resp *types.GetBotInfoRes, err error) {
	// todo: add your logic here and delete this line

	return
}
