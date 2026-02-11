package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT密钥
var jwtSecret = []byte("chprobe_secret_key")

// 自定义Claims
type Claims struct {
	UserUUID string `json:"user_uuid"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 生成JWT令牌
func GenerateToken(userUUID string, username string) (string, error) {
	// 设置过期时间
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建Claims
	claims := &Claims{
		UserUUID: userUUID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "chprobe",
		},
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并返回Token字符串
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析JWT令牌
func ParseToken(tokenString string) (*Claims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证Token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
