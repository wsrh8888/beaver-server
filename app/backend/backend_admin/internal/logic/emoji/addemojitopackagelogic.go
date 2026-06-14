package logic

import (
	"context"
	"errors"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/app/emoji/emoji_rpc/types/emoji_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

const emojiPackageContentActionAdd int32 = 1

type AddEmojiToPackageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddEmojiToPackageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddEmojiToPackageLogic {
	return &AddEmojiToPackageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

// AddEmojiToPackage 管理后台：向表情包添加表情。
// admin 职责：校验运营录入的 packageId/fileKey/title，映射 action=添加。
// RPC 职责：UpdateEmojiPackageContent 统一维护包-表情关联（创建表情 + 建立 relation + 版本号）。
func (l *AddEmojiToPackageLogic) AddEmojiToPackage(req *types.AddEmojiToPackageReq) (resp *types.AddEmojiToPackageRes, err error) {
	if req.PackageId == "" {
		return nil, errors.New("表情包ID不能为空")
	}
	if req.FileUrl == "" {
		return nil, errors.New("文件地址不能为空")
	}
	if req.Title == "" {
		return nil, errors.New("表情名称不能为空")
	}

	rpcRes, err := l.svcCtx.EmojiRpc.UpdateEmojiPackageContent(l.ctx, &emoji_rpc.UpdateEmojiPackageContentReq{
		PackageId: req.PackageId,
		Action:    emojiPackageContentActionAdd,
		FileKey:   req.FileUrl,
		Title:     req.Title,
		EmojiInfo: &emoji_rpc.EmojiInfoMsg{
			Width:  int32(req.EmojiInfo.Width),
			Height: int32(req.EmojiInfo.Height),
		},
	})
	if err != nil {
		l.Errorf("添加表情到表情包失败: %v", err)
		return nil, err
	}
	return &types.AddEmojiToPackageRes{RelationId: rpcRes.RelationId}, nil
}
