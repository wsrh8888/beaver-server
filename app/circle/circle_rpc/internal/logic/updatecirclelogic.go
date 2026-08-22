package logic

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_models"
	"beaver/app/circle/circle_rpc/internal/svc"
	"beaver/app/circle/circle_rpc/types/circle_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCircleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCircleLogic {
	return &UpdateCircleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCircleLogic) UpdateCircle(in *circle_rpc.UpdateCircleReq) (*circle_rpc.UpdateCircleRes, error) {
	updates := map[string]interface{}{}

	if in.IsDeleted != nil {
		updates["is_deleted"] = in.GetIsDeleted()
	}

	if len(updates) == 0 {
		return &circle_rpc.UpdateCircleRes{}, nil
	}

	updates["version"] = l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", in.CircleId)

	if err := l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ?", in.CircleId).
		Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新圈子失败: %v", err)
	}

	return &circle_rpc.UpdateCircleRes{}, nil
}
