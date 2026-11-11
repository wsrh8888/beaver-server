package logic

import (
	"context"
	"sort"

	doc_models "beaver/app/doc/doc_models"
	"beaver/app/doc/doc_rpc/internal/svc"
	"beaver/app/doc/doc_rpc/types/doc_rpc"
	beaverlog "beaver/utils/beaverlog"
)

type ListCollaboratorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *beaverlog.Logger
}

func NewListCollaboratorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCollaboratorsLogic {
	return &ListCollaboratorsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		logger: beaverlog.New("ListCollaboratorsLogic", ctx),
	}
}

func (l *ListCollaboratorsLogic) ListCollaborators(in *doc_rpc.ListCollaboratorsReq) (*doc_rpc.ListCollaboratorsRes, error) {
	doc, _, err := requireDoc(l.svcCtx.DB, in.DocId, in.UserId, roleViewer)
	if err != nil {
		return nil, err
	}
	byUser := map[string]*doc_rpc.CollaboratorItem{}
	if uid := trim(doc.OwnerUserID); uid != "" {
		byUser[uid] = &doc_rpc.CollaboratorItem{
			UserId: uid, Role: roleOwner, PrincipalType: "user", Source: "owner",
		}
	}
	var space doc_models.DocSpace
	if err := l.svcCtx.DB.Where("space_id = ?", doc.SpaceID).First(&space).Error; err == nil {
		if uid := trim(space.OwnerUserID); uid != "" {
			if _, ok := byUser[uid]; !ok {
				byUser[uid] = &doc_rpc.CollaboratorItem{
					UserId: uid, Role: roleOwner, PrincipalType: "user", Source: "space",
				}
			}
		}
	}
	var perms []doc_models.DocPermission
	_ = l.svcCtx.DB.Where(
		"(resource_type = ? AND resource_id = ?) OR (resource_type = ? AND resource_id = ?)",
		"doc", doc.DocID, "space", doc.SpaceID).Find(&perms).Error
	for i := range perms {
		p := &perms[i]
		if p.PrincipalType != "user" || trim(p.PrincipalID) == "" {
			continue
		}
		source := "doc"
		if p.ResourceType == "space" {
			source = "space"
		}
		if cur, ok := byUser[p.PrincipalID]; ok {
			if roleRank(p.Role) > roleRank(cur.Role) {
				cur.Role, cur.Source = p.Role, source
			}
		} else {
			byUser[p.PrincipalID] = &doc_rpc.CollaboratorItem{
				UserId: p.PrincipalID, Role: p.Role, PrincipalType: p.PrincipalType, Source: source,
			}
		}
	}
	out := make([]*doc_rpc.CollaboratorItem, 0, len(byUser))
	for _, v := range byUser {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		if roleRank(out[i].Role) != roleRank(out[j].Role) {
			return roleRank(out[i].Role) > roleRank(out[j].Role)
		}
		return out[i].UserId < out[j].UserId
	})
	return &doc_rpc.ListCollaboratorsRes{Collaborators: out}, nil
}
