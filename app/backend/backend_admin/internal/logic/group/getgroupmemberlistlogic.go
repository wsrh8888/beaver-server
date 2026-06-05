package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/group/group_rpc/types/group_rpc"
	"beaver/app/user/user_rpc/types/user_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *GetGroupMemberListLogic) GetGroupMemberList(req *types.GetGroupMemberListReq) (resp *types.GetGroupMemberListRes, err error) {
	if req.GroupId == "" {
		return nil, errors.New("群组ID不能为空")
	}

	rpcRes, err := l.svcCtx.GroupRpc.ListGroupMembers(l.ctx, &group_rpc.ListGroupMembersReq{
		GroupId:  req.GroupId,
		Page:     int32(req.Page),
		PageSize: int32(req.Limit),
		Role:     int32(req.Role),
		Status:   int32(req.Status),
	})
	if err != nil {
		l.Errorf("获取群成员列表失败: %v", err)
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

	list := make([]types.GetGroupMemberListItem, 0, len(rpcRes.List))
	for _, m := range rpcRes.List {
		nick := ""
		if u, ok := users[m.UserId]; ok && u != nil {
			nick = u.NickName
		}
		list = append(list, types.GetGroupMemberListItem{
			Id:              uint(m.Id),
			GroupId:         m.GroupId,
			UserId:          m.UserId,
			MemberNickname:  nick,
			Role:            int(m.Role),
			ProhibitionTime: int(m.ProhibitionMinutes),
			Status:          int(m.Status),
			CreatedAt:       m.CreatedAt,
			UpdatedAt:       m.UpdatedAt,
		})
	}
	return &types.GetGroupMemberListRes{List: list, Total: rpcRes.Total}, nil
}
