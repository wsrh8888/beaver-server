package logic

import (
	"context"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"

	"beaver/utils/beaverlog"
	"beaver/utils/beaverlog/model"
)

type CreateNodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewCreateNodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateNodeLogic {
	return &CreateNodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("CreateNodeLogic", ctx),
	}
}

func (l *CreateNodeLogic) CreateNode(in *doc_rpc.CreateNodeReq) (*doc_rpc.CreateNodeRes, error) {
	space, err := loadSpace(l.svcCtx.DB, in.SpaceId)
	if err != nil {
		return nil, errNotFound
	}
	role, err := resolveSpaceRole(l.svcCtx.DB, space, in.OperatorUserId)
	if err != nil {
		return nil, err
	}
	if !canEdit(role) {
		return nil, errNoPermission
	}

	nodeType := trim(in.Type)
	if nodeType == "" {
		nodeType = "page"
	}
	title := trim(in.Title)
	if title == "" {
		title = "未命名"
	}

	nodeID := genID("node")
	node := doc_models.DocNode{
		NodeID:   nodeID,
		SpaceID:  space.SpaceID,
		ParentID: trim(in.ParentId),
		Title:    title,
		Icon:     trim(in.Icon),
		Type:     nodeType,
		Sort:     int(in.Sort),
	}

	// page 节点同步创建一篇文档（权威正文源），folder 不建文档。
	var docID string
	if nodeType == "page" {
		docID = genID("doc")
		doc := doc_models.DocDocument{
			DocID:       docID,
			NodeID:      nodeID,
			SpaceID:     space.SpaceID,
			Title:       title,
			ContentJSON: skeletonContentJSON,
			Version:     1,
			OwnerUserID: trim(in.OperatorUserId),
			RagEnabled:  true,
			Status:      1,
		}
		doc.ContentDigest = sha256hex(doc.ContentText)
		node.DocID = docID
		if cerr := l.svcCtx.DB.Create(&doc).Error; cerr != nil {
			l.logger.Error(model.LogMsg{Text: "创建节点文档失败", Data: cerr})
			return nil, cerr
		}
	}

	if err := l.svcCtx.DB.Create(&node).Error; err != nil {
		l.logger.Error(model.LogMsg{Text: "创建节点失败", Data: err})
		return nil, err
	}

	return &doc_rpc.CreateNodeRes{NodeId: nodeID, DocId: docID}, nil
}
