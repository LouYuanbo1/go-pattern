package jwt

import (
	"fmt"
	"go-pattern/internal/config"
	"time"

	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

/*
jwt库中结构

	type RegisteredClaims struct {
		//签发者:通常为 URL 或唯一标识符字符串
		// the `iss` (Issuer) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.1
		Issuer string `json:"iss,omitempty"`

		//主题:通常为用户ID或唯一标识符字符串
		// the `sub` (Subject) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.2
		Subject string `json:"sub,omitempty"`

		//受众:通常为接收 JWT 的应用程序或服务的 URL 或唯一标识符字符串
		// the `aud` (Audience) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.3
		Audience ClaimStrings `json:"aud,omitempty"`

		//过期时间:通常为 JWT 过期的时间戳
		// the `exp` (Expiration Time) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.4
		ExpiresAt *NumericDate `json:"exp,omitempty"`

		//生效时间:通常为 JWT 生效的时间戳
		// the `nbf` (Not Before) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.5
		NotBefore *NumericDate `json:"nbf,omitempty"`

		//签发时间:通常为 JWT 签发的时间戳
		// the `iat` (Issued At) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.6
		IssuedAt *NumericDate `json:"iat,omitempty"`

		//JWT ID:通常为 JWT 的唯一标识符字符串
		// the `jti` (JWT ID) claim. See https://datatracker.ietf.org/doc/html/rfc7519#section-4.1.7
		ID string `json:"jti,omitempty"`
	}
*/
type AuthService struct {
	secretKey   []byte
	tokenExpire int64
	issuer      string
	audience    []string
}

func NewAuthService(jwtConfig config.JWTConfig, issuer string, audience []string) *AuthService {
	return &AuthService{
		secretKey:   []byte(jwtConfig.SecretKey),
		tokenExpire: jwtConfig.TokenExpire,
		issuer:      issuer,
		audience:    audience,
	}
}

type CustomClaims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (a *AuthService) GenerateToken(id uint64, subject string) (string, error) {
	// 生成 JWT 令牌
	claims := CustomClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(a.tokenExpire) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    a.issuer,
			Subject:   subject,
			ID:        uuid.New().String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("GenerateToken: token 生成失败(Token Generate Failed): %w", err)
	}
	return tokenString, nil
}

func (a *AuthService) ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 JWT 令牌
	claims := &CustomClaims{}
	//interface{} 是一个空接口类型，它可以存储任意类型的值。高版本可以换成any
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 关键安全检查：验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("GenerateToken: unexpected signing method: %v", token.Header["alg"])
		}
		return a.secretKey, nil
	})
	if err != nil {
		//作用：检查错误 err 是否与目标错误 target 匹配
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("GenerateToken: token 已过期(Token Expired): %w", err)
		}
		return nil, fmt.Errorf("GenerateToken: token 解析失败(Token Parse Failed): %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("GenerateToken: token 无效(Token Invalid)")
	}
	return claims, nil
}
