package logic

import (
	"context"
	"errors"
	"fmt"

	"beaver/app/chat/chat_rpc/types/chat_rpc"
	"beaver/app/document/document_api/internal/svc"
	"beaver/app/document/document_api/internal/types"
	"beaver/app/document/document_models"
	"beaver/common/models/ctype"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ShareDocumentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewShareDocumentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShareDocumentLogic {
	return &ShareDocumentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ShareDocumentLogic) ShareDocument(req *types.ShareDocumentReq) (*types.ShareDocumentRes, error) {
	var doc document_models.CloudDocument
	err := l.svcCtx.DB.Where("doc_id = ? AND status = 1 AND deleted_at IS NULL", req.DocID).First(&doc).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("文档不存在")
	}
	if err != nil {
		return nil, err
	}

	perm, err := resolveDocumentPerm(l.svcCtx.DB, &doc, req.UserID)
	if err != nil {
		return nil, err
	}
	if !canRead(perm) {
		return nil, fmt.Errorf("无权限分享该文档")
	}

	rpcRes, err := l.svcCtx.ChatRpc.SendMsg(l.ctx, &chat_rpc.SendMsgReq{
		UserId:         req.UserID,
		ConversationId: req.ConversationID,
		MessageId:      req.MessageID,
		Msg: &chat_rpc.Msg{
			Type: uint32(ctype.CloudDocMsgType),
			CloudDocMsg: &chat_rpc.CloudDocMsg{
				DocId:    doc.DocID,
				DocType:  int32(doc.DocType),
				Title:    doc.Title,
				OwnerId:  doc.OwnerID,
				Perm:     int32(perm),
				CoverUrl: doc.CoverURL,
				Revision: doc.Revision,
			},
		},
	})
	if err != nil {
		l.Errorf("分享云文档到会话失败: docId=%s, err=%v", req.DocID, err)
		return nil, fmt.Errorf("发送到会话失败")
	}

	return &types.ShareDocumentRes{
		MessageID: rpcRes.MessageId,
		Seq:       rpcRes.Seq,
	}, nil
}
