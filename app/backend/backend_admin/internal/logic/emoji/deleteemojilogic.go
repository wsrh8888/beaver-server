package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteEmojiLogic {
	return &DeleteEmojiLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// DeleteEmoji 管理后台：删除表情。
// admin 职责：校验路径参数，将「删除」接口语义映射为 SaveEmoji.delete=true。
// RPC 职责：执行领域删除与版本同步，不与 HTTP 路由 1:1 命名。
func (l *DeleteEmojiLogic) DeleteEmoji(req *types.DeleteEmojiReq) (resp *types.DeleteEmojiRes, err error) {
	if req.EmojiId == "" {
		return nil, errors.New("表情ID不能为空")
	}

	del := true
	_, err = l.svcCtx.EmojiRpc.SaveEmoji(l.ctx, &emoji_rpc.SaveEmojiReq{
		EmojiId: req.EmojiId,
		Delete:  &del,
	})
	if err != nil {
		l.Errorf("删除表情失败: %v", err)
		return nil, err
	}
	return &types.DeleteEmojiRes{}, nil
}
