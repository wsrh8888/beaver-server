package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiCollectListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojiCollectListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiCollectListLogic {
	return &GetEmojiCollectListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetEmojiCollectList 管理后台：用户收藏表情列表。
// admin 职责：运营筛选条件映射、响应适配；如需展示收藏者昵称，在此通过 UserRpc.UserListInfo 跨域组装。
// RPC 职责：ListEmojiCollects 查询并 join 表情 title/file_key。
func (l *GetEmojiCollectListLogic) GetEmojiCollectList(req *types.GetEmojiCollectListReq) (resp *types.GetEmojiCollectListRes, err error) {
	rpcRes, err := l.svcCtx.EmojiRpc.ListEmojiCollects(l.ctx, &emoji_rpc.ListEmojiCollectsReq{
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
		UserId:    req.UserID,
		EmojiId:   req.EmojiID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		l.Errorf("获取表情收藏列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetEmojiCollectListItem, 0, len(rpcRes.List))
	for _, c := range rpcRes.List {
		list = append(list, types.GetEmojiCollectListItem{
			CollectId:    c.CollectId,
			UserID:       c.UserId,
			EmojiId:      c.EmojiId,
			EmojiTitle:   c.EmojiTitle,
			EmojiFileUrl: c.EmojiFileKey,
			CreateTime:   c.CreatedAt,
			UpdateTime:   c.UpdatedAt,
		})
	}
	return &types.GetEmojiCollectListRes{List: list, Total: rpcRes.Total}, nil
}
