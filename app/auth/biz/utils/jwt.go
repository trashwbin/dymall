package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	SecretKey = []byte("your-256-bit-secret") // 在生产环境中应该从配置中读取

	ErrTokenExpired     = errors.New("令牌已过期")
	ErrTokenNotValidYet = errors.New("令牌尚未生效")
	ErrTokenMalformed   = errors.New("令牌格式错误")
	ErrTokenInvalid     = errors.New("无效的令牌")
)

// CustomClaims 自定义的JWT声明结构
type CustomClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT令牌
func GenerateToken(userID int64, username, role string) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24小时后过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, ErrTokenMalformed
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, ErrTokenNotValidYet
		} else {
			return nil, ErrTokenInvalid
		}
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
