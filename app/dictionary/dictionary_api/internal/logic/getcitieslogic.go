package logic

import (
	"context"

	"beaver/app/dictionary/dictionary_api/internal/svc"
	"beaver/app/dictionary/dictionary_api/internal/types"
	"beaver/app/dictionary/dictionary_rpc/types/dictionary_rpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCitiesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCitiesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCitiesLogic {
	return &GetCitiesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCitiesLogic) GetCities() (resp *types.GetCitiesRes, err error) {
	rpcResp, err := l.svcCtx.DictionaryRpc.GetCities(l.ctx, &dictionary_rpc.GetCitiesReq{})
	if err != nil {
		return nil, err
	}

	// 转换 RPC 响应为 API 响应
	cities := make([]types.CityInfo, 0, len(rpcResp.Cities))
	for _, city := range rpcResp.Cities {
		cities = append(cities, types.CityInfo{
			Code: city.CityId,
			Name: city.CityName,
		})
	}

	return &types.GetCitiesRes{
		Cities: cities,
	}, nil
}
