package logic

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetEmojiListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetEmojiListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetEmojiListLogic {
	return &GetEmojiListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// GetEmojiList 管理后台：表情列表查询。
// admin 职责：HTTP 入参校验与映射、运营筛选条件组装、响应字段适配前端协议。
// RPC 职责：表情领域数据查询（ListEmojis），不与本接口 1:1。
func (l *GetEmojiListLogic) GetEmojiList(req *types.GetEmojiListReq) (resp *types.GetEmojiListRes, err error) {
	rpcRes, err := l.svcCtx.EmojiRpc.ListEmojis(l.ctx, &emoji_rpc.ListEmojisReq{
		Page:      int32(req.Page),
		PageSize:  int32(req.PageSize),
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		l.Errorf("获取表情列表失败: %v", err)
		return nil, err
	}

	list := make([]types.GetEmojiListItem, 0, len(rpcRes.List))
	for _, e := range rpcRes.List {
		list = append(list, types.GetEmojiListItem{
			EmojiId:    e.EmojiId,
			FileUrl:    e.FileKey,
			Title:      e.Title,
			AuthorID:   req.AuthorID, // 领域表无 author 字段，仅回显筛选条件
			CreateTime: e.CreatedAt,
			UpdateTime: e.UpdatedAt,
		})
	}
	return &types.GetEmojiListRes{List: list, Total: rpcRes.Total}, nil
}
