package logic

import (
	"context"

	"beaver/app/group/group_api/internal/svc"
	"beaver/app/group/group_api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAnnouncementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAnnouncementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAnnouncementLogic {
	return &UpdateAnnouncementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAnnouncementLogic) UpdateAnnouncement(req *types.GroupAnnouncementReq) (resp *types.GroupAnnouncementRes, err error) {
	// todo: add your logic here and delete this line

	return
}
