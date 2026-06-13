package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiPackageEmojisLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojiPackageEmojisLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiPackageEmojisLogic {
	return &GetEmojiPackageEmojisLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetEmojiPackageEmojis 管理后台：查询包内表情列表。
// admin 职责：校验 packageId，复用 ListEmojis 的 package_id 筛选（不与 RPC 新增 1:1 方法）。
// RPC 职责：按包内 sort_order 返回表情数据。
func (l *GetEmojiPackageEmojisLogic) GetEmojiPackageEmojis(req *types.GetEmojiPackageEmojisReq) (resp *types.GetEmojiPackageEmojisRes, err error) {
	if req.PackageId == "" {
		return nil, errors.New("表情包ID不能为空")
	}

	rpcRes, err := l.svcCtx.EmojiRpc.ListEmojis(l.ctx, &emoji_rpc.ListEmojisReq{
		PackageId: req.PackageId,
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
	})
	if err != nil {
		l.Errorf("获取表情包表情列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetEmojiPackageEmojisItem, 0, len(rpcRes.List))
	for _, e := range rpcRes.List {
		list = append(list, types.GetEmojiPackageEmojisItem{
			EmojiId:    e.EmojiId,
			FileUrl:    e.FileKey,
			Title:      e.Title,
			CreateTime: e.CreatedAt,
			UpdateTime: e.UpdatedAt,
		})
	}
	return &types.GetEmojiPackageEmojisRes{List: list, Total: rpcRes.Total}, nil
}
