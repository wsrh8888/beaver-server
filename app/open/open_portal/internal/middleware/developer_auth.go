package middleware

import (
	"context"
	"net/http"
	"strings"

	models "beaver/app/open/open_models"
	"beaver/utils/jwts"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

// DeveloperAuthMiddleware 开发者认证中间件
func DeveloperAuthMiddleware(secretKey string, db *gorm.DB) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 0. 判断是否是公开接口（如登录接口），如果是则直接放行
			if isPublicPath(r.URL.Path) {
				logx.Infof("访问公开接口: path=%s", r.URL.Path)
				next(w, r)
				return
			}

			// 1. 提取 Token
			token := getTokenFromRequest(r)
			if token == "" {
				http.Error(w, `{"code":401,"msg":"缺少认证令牌"}`, http.StatusUnauthorized)
				return
			}

			// 2. 验证 JWT Token
			claims, err := jwts.ParseToken(token, secretKey)
			if err != nil {
				logx.Errorf("JWT 验证失败: %v", err)
				http.Error(w, `{"code":401,"msg":"令牌无效或已过期"}`, http.StatusUnauthorized)
				return
			}

			// 3. 注入用户ID到 context (供后续 handler 使用)
			ctx := context.WithValue(r.Context(), "userId", claims.UserID)
			r = r.WithContext(ctx)
			r.Header.Set("Beaver-User-Id", claims.UserID)

			// 4. 判断是否是申请开发者接口 (不需要检查开发者资质，但仍需验证 Token)
			if isApplyDeveloperPath(r.URL.Path) {
				// 申请接口只需要登录,不检查开发者资质
				logx.Infof("用户访问申请接口: user_id=%s", claims.UserID)
				next(w, r)
				return
			}

			// 5. 判断是否是只读接口 (查看列表、详情等，不需要开发者资质)
			if isReadOnlyPath(r.URL.Path) {
				// 只读接口只需要登录，不检查开发者资质
				logx.Infof("用户访问只读接口: user_id=%s, path=%s", claims.UserID, r.URL.Path)
				next(w, r)
				return
			}

			// 6. 其他接口(创建、更新、删除等)需要检查开发者资质
			var developer models.OpenDeveloper
			err = db.Where("user_id = ? AND status = ?", claims.UserID, 1).First(&developer).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					http.Error(w, `{"code":403,"msg":"您还不是认证开发者,请先申请开发者资质"}`, http.StatusForbidden)
					return
				}
				logx.Errorf("查询开发者信息失败: %v", err)
				http.Error(w, `{"code":500,"msg":"服务内部异常"}`, http.StatusInternalServerError)
				return
			}

			// 7. 注入开发者ID到 header
			r.Header.Set("Beaver-Developer-Id", string(developer.Id))

			logx.Infof("开发者认证成功: user_id=%s, developer_id=%d", claims.UserID, developer.Id)

			next(w, r)
		}
	}
}

// getTokenFromRequest 从请求中获取 Token
func getTokenFromRequest(r *http.Request) string {
	// 优先从 Authorization header 获取
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// 其次从 Token header 获取
	token := r.Header.Get("Token")
	if token != "" {
		return token
	}

	// 最后从 query 参数获取
	return r.URL.Query().Get("token")
}

// isApplyDeveloperPath 判断是否是申请开发者接口
func isApplyDeveloperPath(path string) bool {
	return strings.HasPrefix(path, "/portal/open/v1/developer/apply")
}

// isPublicPath 判断是否是公开接口（不需要认证）
func isPublicPath(path string) bool {
	publicPaths := []string{
		"/portal/open/v1/auth/login", // 登录接口
	}

	for _, p := range publicPaths {
		if path == p || strings.HasPrefix(path, p+"?") {
			return true
		}
	}

	return false
}

// isReadOnlyPath 判断是否是只读接口 (不需要开发者资质)
func isReadOnlyPath(path string) bool {
	// 应用列表、详情等只读接口
	readOnlyPaths := []string{
		"/portal/open/v1/app/list",     // 应用列表
		"/portal/open/v1/app/detail",   // 应用详情
		"/portal/open/v1/webhook/logs", // Webhook 日志
	}

	for _, p := range readOnlyPaths {
		if path == p || strings.HasPrefix(path, p+"?") {
			return true
		}
	}

	return false
}
