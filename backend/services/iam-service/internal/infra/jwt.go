package infra

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenRevocationStore 跨实例共享的 refresh 吊销存储（如 Redis）。
type TokenRevocationStore interface {
	Revoke(ctx context.Context, jti string, ttl time.Duration) error
	IsRevoked(ctx context.Context, jti string) (bool, error)
}

// JWTTokenManager JWT 令牌管理器
type JWTTokenManager struct {
	secret        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	issuer        string
	revoked       sync.Map // local fallback: jti -> expiry unix
	revokeStore   TokenRevocationStore
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

// WithRevocationStore 注入共享吊销存储；未注入时回退进程内 sync.Map。
func (m *JWTTokenManager) WithRevocationStore(store TokenRevocationStore) *JWTTokenManager {
	m.revokeStore = store
	return m
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
			ID:        uuid.New().String(),
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
			ID:        uuid.New().String(),
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

// RevokeRefreshToken 将 refresh token 的 jti 加入吊销表直至过期
func (m *JWTTokenManager) RevokeRefreshToken(tokenString string) error {
	claims, err := m.parseClaims(tokenString)
	if err != nil {
		return err
	}
	if claims.Type != "refresh" {
		return fmt.Errorf("令牌类型不匹配")
	}
	if claims.ID == "" {
		return fmt.Errorf("令牌缺少 jti")
	}
	ttl := m.refreshExpiry
	if claims.ExpiresAt != nil {
		ttl = time.Until(claims.ExpiresAt.Time)
		if ttl <= 0 {
			return nil
		}
	}
	if m.revokeStore != nil {
		if err := m.revokeStore.Revoke(context.Background(), claims.ID, ttl); err != nil {
			// fall back to local map if redis hiccups
			m.revoked.Store(claims.ID, time.Now().Add(ttl).Unix())
			return nil
		}
		return nil
	}
	m.revoked.Store(claims.ID, time.Now().Add(ttl).Unix())
	return nil
}

func (m *JWTTokenManager) parseClaims(tokenString string) (*jwtClaims, error) {
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
	return claims, nil
}

func (m *JWTTokenManager) isRevoked(jti string) bool {
	if jti == "" {
		return false
	}
	if m.revokeStore != nil {
		ok, err := m.revokeStore.IsRevoked(context.Background(), jti)
		if err == nil {
			return ok
		}
	}
	if expRaw, ok := m.revoked.Load(jti); ok {
		exp, _ := expRaw.(int64)
		if time.Now().Unix() < exp {
			return true
		}
		m.revoked.Delete(jti)
	}
	return false
}

func (m *JWTTokenManager) validateToken(tokenString, expectedType string) (*domain.TokenClaims, error) {
	claims, err := m.parseClaims(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != expectedType {
		return nil, fmt.Errorf("令牌类型不匹配")
	}

	if expectedType == "refresh" && m.isRevoked(claims.ID) {
		return nil, fmt.Errorf("令牌已吊销")
	}

	return &domain.TokenClaims{
		UserID:   claims.UserID,
		TenantID: claims.TenantID,
		Username: claims.Username,
		Roles:    claims.Roles,
	}, nil
}
