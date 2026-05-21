package infra

import (
	"fmt"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// JWTTokenManager JWT 令牌管理器
type JWTTokenManager struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
}

// NewJWTTokenManager 创建 JWT 令牌管理器
func NewJWTTokenManager(secret string, accessExpiry, refreshExpiry time.Duration, issuer string) *JWTTokenManager {
	return &JWTTokenManager{
		secret:        secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		issuer:        issuer,
	}
}

type jwtClaims struct {
	UserID   string   `json:"user_id"`
	TenantID string   `json:"tenant_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Type     string   `json:"type"`
	jwt.RegisteredClaims
}

func (m *JWTTokenManager) GenerateAccessToken(userID, tenantID string, roles []domain.Role) (string, error) {
	roleCodes := make([]string, len(roles))
	for i, r := range roles {
		roleCodes[i] = r.Code
	}

	claims := jwtClaims{
		UserID:   userID,
		TenantID: tenantID,
		Roles:    roleCodes,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTTokenManager) GenerateRefreshToken(userID, tenantID string) (string, error) {
	claims := jwtClaims{
		UserID:   userID,
		TenantID: tenantID,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secret))
}

func (m *JWTTokenManager) ValidateAccessToken(tokenString string) (*domain.TokenClaims, error) {
	return m.validateToken(tokenString, "access")
}

func (m *JWTTokenManager) ValidateRefreshToken(tokenString string) (*domain.TokenClaims, error) {
	return m.validateToken(tokenString, "refresh")
}

func (m *JWTTokenManager) validateToken(tokenString, expectedType string) (*domain.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return []byte(m.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("无效的令牌")
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("令牌类型不匹配")
	}

	return &domain.TokenClaims{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Username: claims.Username,
		Roles:    claims.Roles,
	}, nil
}
