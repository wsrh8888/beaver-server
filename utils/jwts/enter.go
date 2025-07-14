package jwts

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtPayLoad struct {
	Nickname string `json:"nickname"`
	UserID   string `json:"userId"`
}

type CustomClaims struct {
	JwtPayLoad
	jwt.RegisteredClaims
}

func GenToken(payload JwtPayLoad, accessSecret string, expires int) (string, error) {
	claims := CustomClaims{
		JwtPayLoad: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expires))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessSecret))
}

func ParseToken(tokenStr string, accessSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if clains, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return clains, nil
	}
	return nil, err
}
