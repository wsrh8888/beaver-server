package circle

import (
	"context"
	"fmt"

	"beaver/app/circle/circle_api/internal/svc"
	"beaver/app/circle/circle_api/internal/types"
	"beaver/app/circle/circle_models"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuitCircleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuitCircleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuitCircleLogic {
	return &QuitCircleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuitCircleLogic) QuitCircle(req *types.QuitCircleReq) (resp *types.QuitCircleRes, err error) {
	var member circle_models.CircleMemberModel
	if err = l.svcCtx.DB.Where("circle_id = ? AND user_id = ?", req.CircleID, req.UserID).First(&member).Error; err != nil {
		return nil, fmt.Errorf("你不是该圈子成员")
	}
	if member.Role == 1 {
		return nil, fmt.Errorf("圈主不能退出，请先转让圈主")
	}

	if err = l.svcCtx.DB.Delete(&member).Error; err != nil {
		return nil, fmt.Errorf("退出圈子失败: %v", err)
	}

	// 更新圈子版本和成员数
	circleVersion := l.svcCtx.VersionGen.GetNextVersion("circles", "circle_id", req.CircleID)
	l.svcCtx.DB.Model(&circle_models.CircleModel{}).
		Where("circle_id = ? AND member_count > 0", req.CircleID).
		Updates(map[string]interface{}{
			"member_count": l.svcCtx.DB.Raw("GREATEST(member_count - 1, 0)"),
			"version":      circleVersion,
		})

	return &types.QuitCircleRes{}, nil
}
