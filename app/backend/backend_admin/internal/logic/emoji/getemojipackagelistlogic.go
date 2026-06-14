package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojiPackageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageListLogic {
	return &GetEmojiPackageListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetEmojiPackageList 管理后台：表情包列表查询。
// admin 职责：运营筛选条件映射、响应适配；如需展示创建者昵称，在此通过 UserRpc.UserListInfo 跨域组装。
// RPC 职责：ListEmojiPackages 领域查询。
func (l *GetEmojiPackageListLogic) GetEmojiPackageList(req *types.GetEmojiPackageListReq) (resp *types.GetEmojiPackageListRes, err error) {
	rpcRes, err := l.svcCtx.EmojiRpc.ListEmojiPackages(l.ctx, &emoji_rpc.ListEmojiPackagesReq{
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
		UserId:    req.UserID,
		Type:      req.Type,
		Status:    int32(req.Status),
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		l.Errorf("获取表情包列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetEmojiPackageListItem, 0, len(rpcRes.List))
	for _, p := range rpcRes.List {
		list = append(list, types.GetEmojiPackageListItem{
			PackageId:   p.PackageId,
			Title:       p.Title,
			CoverFile:   p.CoverFile,
			UserID:      p.UserId,
			Description: p.Description,
			Type:        p.Type,
			Status:      int(p.Status),
			CreateTime:  p.CreatedAt,
			UpdateTime:  p.UpdatedAt,
		})
	}
	return &types.GetEmojiPackageListRes{List: list, Total: rpcRes.Total}, nil
}
