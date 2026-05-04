package group

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取群组信息
func NewGetGroupInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupInfoLogic {
	return &GetGroupInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGroupInfoLogic) GetGroupInfo(req *types.GetGroupInfoReq) (resp *types.GetGroupInfoRes, err error) {
	// TODO: 需要调用 Group RPC 获取群组信息
	logx.Infof("获取群组信息: groupID=%s", req.GroupID)

	return &types.GetGroupInfoRes{
		Group: types.GroupInfo{
			GroupID: req.GroupID,
			Name:    "示例群组",
		},
	}, nil
}
