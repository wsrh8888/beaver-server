package group

import (
	"context"

	"beaver/app/open/open_api/internal/svc"
	"beaver/app/open/open_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGroupMemberLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加群成员
func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddGroupMemberLogic) AddGroupMember(req *types.AddGroupMemberReq) (resp *types.AddGroupMemberRes, err error) {
	// TODO: 需要调用 Group RPC 添加群成员
	logx.Infof("添加群成员: groupID=%s, memberCount=%d", req.GroupID, len(req.MemberIDs))

	return &types.AddGroupMemberRes{
		Success: true,
	}, nil
}
