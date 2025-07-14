package handler

import (
	"beaver/app/dictionary/dictionary_api/internal/logic"
	"beaver/app/dictionary/dictionary_api/internal/svc"
	"beaver/common/response"
	"net/http"
)

func GetCitiesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewGetCitiesLogic(r.Context(), svcCtx)
		resp, err := l.GetCities()
		response.Response(r, w, resp, err)
	}
}
