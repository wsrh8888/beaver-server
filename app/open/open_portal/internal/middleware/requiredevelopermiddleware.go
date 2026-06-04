package middleware

import (
	"net/http"

	models "beaver/app/open/open_models"

	"gorm.io/gorm"
)

type RequireDeveloperMiddleware struct {
	db *gorm.DB
}

func NewRequireDeveloperMiddleware(db *gorm.DB) *RequireDeveloperMiddleware {
	return &RequireDeveloperMiddleware{db: db}
}

func (m *RequireDeveloperMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Beaver-User-Id")
		if userID == "" {
			http.Error(w, `{"code":401,"msg":"未登录"}`, http.StatusUnauthorized)
			return
		}

		var developer models.OpenDeveloper
		err := m.db.Where("user_id = ? AND status = ?", userID, 1).First(&developer).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				http.Error(w, `{"code":403,"msg":"您还不是认证开发者,请先申请开发者资质"}`, http.StatusForbidden)
				return
			}
			http.Error(w, `{"code":500,"msg":"服务内部异常"}`, http.StatusInternalServerError)
			return
		}

		next(w, r)
	}
}
