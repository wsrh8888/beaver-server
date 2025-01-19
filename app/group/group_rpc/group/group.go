// Code generated by goctl. DO NOT EDIT.
// Source: group_rpc.proto

package group

import (
	"context"

	"beaver/app/group/group_rpc/types/group_rpc"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GetGroupMembersReq = group_rpc.GetGroupMembersReq
	GetGroupMembersRes = group_rpc.GetGroupMembersRes
	GroupMemberInfo    = group_rpc.GroupMemberInfo

	Group interface {
		GetGroupMembers(ctx context.Context, in *GetGroupMembersReq, opts ...grpc.CallOption) (*GetGroupMembersRes, error)
	}

	defaultGroup struct {
		cli zrpc.Client
	}
)

func NewGroup(cli zrpc.Client) Group {
	return &defaultGroup{
		cli: cli,
	}
}

func (m *defaultGroup) GetGroupMembers(ctx context.Context, in *GetGroupMembersReq, opts ...grpc.CallOption) (*GetGroupMembersRes, error) {
	client := group_rpc.NewGroupClient(m.cli.Conn())
	return client.GetGroupMembers(ctx, in, opts...)
}
