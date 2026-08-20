package app

import (
	"context"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/stretchr/testify/require"
)

type memSession struct {
	data map[string]string
}

func (m *memSession) Put(_ context.Context, id, val string, _ time.Duration) error {
	if m.data == nil {
		m.data = map[string]string{}
	}
	m.data[id] = val
	return nil
}

func (m *memSession) Delete(_ context.Context, id string) error {
	delete(m.data, id)
	return nil
}

type stubUserRepo struct {
	user *domain.User
}

func (r *stubUserRepo) Create(context.Context, *domain.User) error { return nil }
func (r *stubUserRepo) Update(context.Context, *domain.User) error  { return nil }
func (r *stubUserRepo) FindByID(context.Context, string, string) (*domain.User, error) {
	return r.user, nil
}
func (r *stubUserRepo) FindByUsername(context.Context, string, string) (*domain.User, error) {
	return r.user, nil
}
func (r *stubUserRepo) FindWithRoles(context.Context, string, string) (*domain.User, error) {
	return r.user, nil
}
func (r *stubUserRepo) List(context.Context, string, int, int) ([]*domain.User, int64, error) {
	return nil, 0, nil
}
func (r *stubUserRepo) Delete(context.Context, string, string) error { return nil }

type stubRoleRepo struct{}

func (stubRoleRepo) Create(context.Context, *domain.Role) error { return nil }
func (stubRoleRepo) Update(context.Context, *domain.Role) error { return nil }
func (stubRoleRepo) FindByID(context.Context, string, string) (*domain.Role, error) {
	return nil, nil
}
func (stubRoleRepo) FindByCode(context.Context, string, string) (*domain.Role, error) {
	return nil, nil
}
func (stubRoleRepo) FindWithPermissions(context.Context, string, string) (*domain.Role, error) {
	return nil, nil
}
func (stubRoleRepo) List(context.Context, string, int, int) ([]*domain.Role, int64, error) {
	return nil, 0, nil
}
func (stubRoleRepo) Delete(context.Context, string, string) error { return nil }
func (stubRoleRepo) AddPermissions(context.Context, string, string, []string) error {
	return nil
}
func (stubRoleRepo) RemovePermissions(context.Context, string, string, []string) error {
	return nil
}
func (stubRoleRepo) AssignUserRoles(context.Context, string, string, []string) error {
	return nil
}
func (stubRoleRepo) RemoveUserRoles(context.Context, string, string, []string) error {
	return nil
}

type stubTokenMgr struct{}

func (stubTokenMgr) GenerateAccessToken(userID, _ string, _ []domain.Role) (string, error) {
	return "access-" + userID, nil
}
func (stubTokenMgr) GenerateRefreshToken(userID, _ string) (string, error) {
	return "refresh-" + userID, nil
}
func (stubTokenMgr) ValidateAccessToken(string) (*domain.TokenClaims, error) { return nil, nil }
func (stubTokenMgr) ValidateRefreshToken(string) (*domain.TokenClaims, error) {
	return &domain.TokenClaims{UserID: "u1", TenantID: "t1"}, nil
}
func (stubTokenMgr) RevokeRefreshToken(string) error { return nil }

type stubHasher struct{}

func (stubHasher) Hash(string) (string, error) { return "hash", nil }
func (stubHasher) Verify(password, hash string) bool {
	return password == "ok" && hash == "hash"
}

type stubAudit struct{}

func (stubAudit) Write(context.Context, *domain.AuditLog) error { return nil }
func (stubAudit) List(context.Context, string, int, int) ([]*domain.AuditLog, int64, error) {
	return nil, 0, nil
}

func TestLoginLogout_WritesSession(t *testing.T) {
	user := &domain.User{
		ID: "u1", TenantID: "t1", Username: "admin",
		PasswordHash: "hash", Status: domain.UserStatusActive,
	}
	sess := &memSession{}
	svc := NewAuthService(&stubUserRepo{user: user}, stubRoleRepo{}, stubTokenMgr{}, stubHasher{}, stubAudit{}).
		WithSessionStore(sess)

	res, err := svc.Login(context.Background(), "t1", "admin", "ok", "127.0.0.1", "test")
	require.NoError(t, err)
	require.NotEmpty(t, res.AccessToken)
	require.Contains(t, sess.data["u1"], `"user_id":"u1"`)

	require.NoError(t, svc.Logout(context.Background(), "t1", "u1", "admin", "127.0.0.1", "test", res.RefreshToken))
	_, ok := sess.data["u1"]
	require.False(t, ok)
}
