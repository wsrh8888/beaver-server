package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateEmojiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateEmojiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateEmojiLogic {
	return &UpdateEmojiLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// UpdateEmoji 管理后台：更新表情。
// admin 职责：校验路径参数 emojiId，将可选更新字段组装为 patch 语义（只传有值的字段）。
// RPC 职责：SaveEmoji 更新分支，处理版本递增与同名冲突等领域规则。
func (l *UpdateEmojiLogic) UpdateEmoji(req *types.UpdateEmojiReq) (resp *types.UpdateEmojiRes, err error) {
	if req.EmojiId == "" {
		return nil, errors.New("表情ID不能为空")
	}

	rpcReq := &emoji_rpc.SaveEmojiReq{EmojiId: req.EmojiId}
	if req.FileKey != nil {
		rpcReq.PatchFileKey = req.FileKey
	}
	if req.Title != nil {
		rpcReq.PatchTitle = req.Title
	}

	_, err = l.svcCtx.EmojiRpc.SaveEmoji(l.ctx, rpcReq)
	if err != nil {
		l.Errorf("更新表情失败: %v", err)
		return nil, err
	}
	return &types.UpdateEmojiRes{}, nil
}
