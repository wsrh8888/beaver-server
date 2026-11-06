package logic

import (
	"context"
	"errors"

	doc_rpc "beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/app/doc/doc_api/internal/svc"
)

// errNoPermission API 层权限预检失败（网关已做 JWT 鉴权，此处做资源级角色校验）。
var errNoPermission = errors.New("doc: 无权限")

func roleRank(role string) int {
	switch role {
	case "owner":
		return 5
	case "manager":
		return 4
	case "editor":
		return 3
	case "commenter":
		return 2
	case "viewer":
		return 1
	default:
		return 0
	}
}

// requireRole 通过 doc_rpc.CheckPermission 校验主体对资源的有效角色是否达到 min 级别。
// 用于 RPC 报文未携带 operator 的写操作（如 UpdateSpace），由 API 层补齐权限判定。
func requireRole(ctx context.Context, svcCtx *svc.ServiceContext, resourceType, resourceID, userID, minRole string) error {
	res, err := svcCtx.Doc.CheckPermission(ctx, &doc_rpc.CheckPermissionReq{
		ResourceType: resourceType,
		ResourceId:   resourceID,
		UserId:       userID,
	})
	if err != nil {
		return err
	}
	if roleRank(res.GetRole()) < roleRank(minRole) {
		return errNoPermission
	}
	return nil
}
