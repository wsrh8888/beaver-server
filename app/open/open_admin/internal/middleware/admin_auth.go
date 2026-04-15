package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
)

// AdminAuthMiddleware 管理后台 JWT 认证中间件
func AdminAuthMiddleware(secretKey string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 1. 提取 Token
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"code":401,"msg":"缺少认证令牌"}`, http.StatusUnauthorized)
				return
			}

			// 2. 验证 Bearer 格式
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"code":401,"msg":"无效的认证格式"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// 3. 验证 JWT
			claims, err := validateJWT(tokenString, secretKey)
			if err != nil {
				logx.Errorf("JWT 验证失败: %v", err)
				http.Error(w, `{"code":401,"msg":"令牌无效或已过期"}`, http.StatusUnauthorized)
				return
			}

			// 4. 检查权限（RBAC）
			if !hasPermission(claims.Role, r.Method, r.URL.Path) {
				http.Error(w, `{"code":403,"msg":"权限不足"}`, http.StatusForbidden)
				return
			}

			// 5. 注入用户信息到 header
			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Role", claims.Role)

			next(w, r)
		}
	}
}

// Claims JWT Claims
type Claims struct {
	UserID string `json:"userId"`
	Role   string `json:"role"` // admin, operator, viewer
	jwt.RegisteredClaims
}

// validateJWT 验证 JWT Token
func validateJWT(tokenString string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// hasPermission RBAC 权限检查
func hasPermission(role, method, path string) bool {
	// 权限矩阵
	permissions := map[string]map[string][]string{
		"admin": {
			"GET":    {"/admin/open/*"},
			"POST":   {"/admin/open/*"},
			"PUT":    {"/admin/open/*"},
			"DELETE": {"/admin/open/*"},
		},
		"operator": {
			"GET":  {"/admin/open/*"},
			"POST": {"/admin/open/app/*", "/admin/open/webhook/*"},
		},
		"viewer": {
			"GET": {"/admin/open/*"},
		},
	}

	allowedPaths, ok := permissions[role][method]
	if !ok {
		return false
	}

	// 路径匹配
	for _, allowedPath := range allowedPaths {
		if matchPath(path, allowedPath) {
			return true
		}
	}

	return false
}

// matchPath 路径匹配（支持通配符 *）
func matchPath(path, pattern string) bool {
	if pattern == "/admin/open/*" {
		return strings.HasPrefix(path, "/admin/open/")
	}
	return path == pattern
}

// GenerateJWT 生成 JWT Token
func GenerateJWT(userID, role, secretKey string, expireHours int) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// 错误定义
var (
	ErrTokenInvalid = &AuthError{Code: 401, Msg: "令牌无效"}
	ErrTokenExpired = &AuthError{Code: 401, Msg: "令牌已过期"}
)

type AuthError struct {
	Code int
	Msg  string
}

func (e *AuthError) Error() string {
	return e.Msg
}
