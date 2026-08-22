package circle

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/circle/circle_rpc/types/circle_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCircleMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCircleMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCircleMemberListLogic {
	return &GetCircleMemberListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCircleMemberListLogic) GetCircleMemberList(req *types.GetCircleMemberListReq) (resp *types.GetCircleMemberListRes, err error) {
	if req.CircleId == "" {
		return nil, errors.New("圈子ID不能为空")
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 50
	}

	rpcRes, err := l.svcCtx.CircleRpc.GetCircleMembers(l.ctx, &circle_rpc.GetCircleMembersReq{
		CircleId: req.CircleId,
		Page:     int32(page),
		PageSize: int32(limit),
	})
	if err != nil {
		l.Errorf("获取圈子成员列表失败: %v", err)
		return nil, err
	}

	users := map[string]*user_rpc.UserInfo{}
	seen := map[string]struct{}{}
	userIDs := make([]string, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		if m.UserId == "" {
			continue
		}
		if _, ok := seen[m.UserId]; ok {
			continue
		}
		seen[m.UserId] = struct{}{}
		userIDs = append(userIDs, m.UserId)
	}
	if len(userIDs) > 0 {
		if res, err := l.svcCtx.UserRpc.UserListInfo(l.ctx, &user_rpc.UserListInfoReq{UserIdList: userIDs}); err == nil && res != nil {
			users = res.UserInfo
		}
	}

	list := make([]types.GetCircleMemberListItem, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		nick := m.UserId
		if u, ok := users[m.UserId]; ok && u != nil && u.NickName != "" {
			nick = u.NickName
		}
		list = append(list, types.GetCircleMemberListItem{
			CircleId: m.CircleId,
			UserId:   m.UserId,
			NickName: nick,
			Role:     int(m.Role),
		})
	}
	return &types.GetCircleMemberListRes{List: list, Total: rpcRes.Total}, nil
}
