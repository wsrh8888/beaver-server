package monitor

import (
	"context"

	"beaver/app/backend/backend_admin/internal/svc"
	"beaver/app/backend/backend_admin/internal/types"
	"beaver/core/coreonline"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOnlineStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 在线用户统计
func NewGetOnlineStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOnlineStatsLogic {
	return &GetOnlineStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOnlineStatsLogic) GetOnlineStats(req *types.GetOnlineStatsReq) (resp *types.GetOnlineStatsRes, err error) {
	online, err := coreonline.List(l.svcCtx.Redis)
	if err != nil {
		l.Errorf("获取在线用户统计失败: %v", err)
		return nil, err
	}

	resp = &types.GetOnlineStatsRes{
		UserCount: int64(len(online)),
	}

	for _, user := range online {
		for _, slot := range user.Slots {
			switch slot.Slot {
			case "desktop":
				resp.DesktopCount++
			case "mobile":
				resp.MobileCount++
			}
		}
	}

	return resp, nil
}
