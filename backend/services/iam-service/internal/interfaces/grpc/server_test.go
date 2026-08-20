package grpciface_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/Tangyd893/ERP-Go/backend/proto/iam/v1"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/infra"
	grpciface "github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/interfaces/grpc"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

type mockUserRepo struct {
	users map[string]*domain.User
}

func (m *mockUserRepo) Create(_ context.Context, user *domain.User) error { return nil }
func (m *mockUserRepo) Update(_ context.Context, user *domain.User) error { return nil }
func (m *mockUserRepo) FindByID(_ context.Context, _, userID string) (*domain.User, error) {
	if u, ok := m.users[userID]; ok {
		return u, nil
	}
	return nil, nil
}
func (m *mockUserRepo) FindByUsername(_ context.Context, _, _ string) (*domain.User, error) {
	return nil, nil
}
func (m *mockUserRepo) FindWithRoles(_ context.Context, _, userID string) (*domain.User, error) {
	return m.FindByID(nil, "", userID)
}
func (m *mockUserRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.User, int64, error) {
	return nil, 0, nil
}
func (m *mockUserRepo) Delete(_ context.Context, _, _ string) error { return nil }

type mockRoleRepo struct{}

func (m *mockRoleRepo) Create(_ context.Context, _ *domain.Role) error { return nil }
func (m *mockRoleRepo) Update(_ context.Context, _ *domain.Role) error { return nil }
func (m *mockRoleRepo) FindByID(_ context.Context, _, _ string) (*domain.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindByCode(_ context.Context, _, _ string) (*domain.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) FindWithPermissions(_ context.Context, _, _ string) (*domain.Role, error) {
	return nil, nil
}
func (m *mockRoleRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.Role, int64, error) {
	return nil, 0, nil
}
func (m *mockRoleRepo) Delete(_ context.Context, _, _ string) error { return nil }
func (m *mockRoleRepo) AddPermissions(_ context.Context, _, _ string, _ []string) error {
	return nil
}
func (m *mockRoleRepo) RemovePermissions(_ context.Context, _, _ string, _ []string) error {
	return nil
}
func (m *mockRoleRepo) AssignUserRoles(_ context.Context, _, _ string, _ []string) error {
	return nil
}
func (m *mockRoleRepo) RemoveUserRoles(_ context.Context, _, _ string, _ []string) error {
	return nil
}

type mockAuditRepo struct{}

func (m *mockAuditRepo) Write(_ context.Context, _ *domain.AuditLog) error { return nil }
func (m *mockAuditRepo) List(_ context.Context, _ string, _, _ int) ([]*domain.AuditLog, int64, error) {
	return nil, 0, nil
}

func setupGRPCClient(t *testing.T) (iamv1.IAMServiceClient, func()) {
	t.Helper()

	userRepo := &mockUserRepo{users: map[string]*domain.User{
		"user-1": {
			ID:       "user-1",
			TenantID: "default",
			Username: "admin",
			Status:   domain.UserStatusActive,
			Roles: []domain.Role{{
				Code: "admin",
				Permissions: []domain.Permission{
					{Code: "order:read"},
				},
			}},
		},
	}}
	hasher := infra.NewBcryptPasswordHasher(10)
	tokenMgr := infra.NewJWTTokenManager("grpc-test-secret", time.Hour, 24*time.Hour, "erp-go")
	authSvc := app.NewAuthService(userRepo, &mockRoleRepo{}, tokenMgr, hasher, &mockAuditRepo{})

	listener := bufconn.Listen(bufSize)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpciface.InternalTokenUnaryServerInterceptor()))
	iamv1.RegisterIAMServiceServer(grpcServer, grpciface.NewServer(authSvc))
	go func() { _ = grpcServer.Serve(listener) }()

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}

	cleanup := func() {
		_ = conn.Close()
		grpcServer.Stop()
	}
	return iamv1.NewIAMServiceClient(conn), cleanup
}

func TestCheckPermission_GRPC_Allowed(t *testing.T) {
	client, cleanup := setupGRPCClient(t)
	defer cleanup()

	resp, err := client.CheckPermission(context.Background(), &iamv1.CheckPermissionRequest{
		UserId:         "user-1",
		TenantId:       "default",
		PermissionCode: "order:read",
	})
	if err != nil {
		t.Fatalf("CheckPermission: %v", err)
	}
	if !resp.GetAllowed() {
		t.Fatalf("期望 allowed=true，实际 false，message=%q", resp.GetMessage())
	}
}

func TestCheckPermission_GRPC_Denied(t *testing.T) {
	client, cleanup := setupGRPCClient(t)
	defer cleanup()

	resp, err := client.CheckPermission(context.Background(), &iamv1.CheckPermissionRequest{
		UserId:         "user-1",
		TenantId:       "default",
		PermissionCode: "order:delete",
	})
	if err != nil {
		t.Fatalf("CheckPermission: %v", err)
	}
	if resp.GetAllowed() {
		t.Fatal("期望 allowed=false")
	}
	if resp.GetMessage() == "" {
		t.Fatal("期望拒绝时返回 message")
	}
}

func TestCheckPermission_GRPC_InvalidArgument(t *testing.T) {
	client, cleanup := setupGRPCClient(t)
	defer cleanup()

	_, err := client.CheckPermission(context.Background(), &iamv1.CheckPermissionRequest{
		UserId:   "user-1",
		TenantId: "default",
	})
	if err == nil {
		t.Fatal("期望 InvalidArgument 错误")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("期望 InvalidArgument，实际 %v", status.Code(err))
	}
}

func TestCheckPermission_GRPC_RequiresInternalToken(t *testing.T) {
	t.Setenv("INTERNAL_SERVICE_TOKEN", "test-internal-token")
	client, cleanup := setupGRPCClient(t)
	defer cleanup()

	_, err := client.CheckPermission(context.Background(), &iamv1.CheckPermissionRequest{
		UserId:         "user-1",
		TenantId:       "default",
		PermissionCode: "order:read",
	})
	if err == nil {
		t.Fatal("期望未携带内部令牌时 Unauthenticated")
	}
	if status.Code(err) != codes.Unauthenticated {
		t.Fatalf("期望 Unauthenticated，实际 %v", status.Code(err))
	}

	md := metadata.Pairs("x-erp-internal-token", "test-internal-token")
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.CheckPermission(ctx, &iamv1.CheckPermissionRequest{
		UserId:         "user-1",
		TenantId:       "default",
		PermissionCode: "order:read",
	})
	if err != nil {
		t.Fatalf("携带令牌后 CheckPermission: %v", err)
	}
	if !resp.GetAllowed() {
		t.Fatal("期望 allowed=true")
	}
}

func TestTryStart_ListensAndStops(t *testing.T) {
	userRepo := &mockUserRepo{users: map[string]*domain.User{}}
	hasher := infra.NewBcryptPasswordHasher(10)
	tokenMgr := infra.NewJWTTokenManager("grpc-test-secret", time.Hour, 24*time.Hour, "erp-go")
	authSvc := app.NewAuthService(userRepo, &mockRoleRepo{}, tokenMgr, hasher, &mockAuditRepo{})
	log := logger.New("info", "text", "stdout", "iam-grpc-test", "testing")

	t.Setenv("IAM_GRPC_ADDR", "127.0.0.1:0")
	stop := grpciface.TryStart(authSvc, log)
	defer stop()

	// 短暂等待 goroutine 完成 bind
	time.Sleep(100 * time.Millisecond)
}

func TestTryStart_Binds9091(t *testing.T) {
	userRepo := &mockUserRepo{users: map[string]*domain.User{}}
	hasher := infra.NewBcryptPasswordHasher(10)
	tokenMgr := infra.NewJWTTokenManager("grpc-test-secret", time.Hour, 24*time.Hour, "erp-go")
	authSvc := app.NewAuthService(userRepo, &mockRoleRepo{}, tokenMgr, hasher, &mockAuditRepo{})
	log := logger.New("info", "text", "stdout", "iam-grpc-test", "testing")

	t.Setenv("IAM_GRPC_ADDR", ":9091")
	stop := grpciface.TryStart(authSvc, log)
	defer stop()

	time.Sleep(200 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", "127.0.0.1:9091", time.Second)
	if err != nil {
		t.Fatalf("无法连接 :9091: %v", err)
	}
	_ = conn.Close()
}
