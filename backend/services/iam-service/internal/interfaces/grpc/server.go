package grpciface

import (
	"context"
	"errors"
	"net"
	"os"
	"strings"

	iamv1 "github.com/Tangyd893/ERP-Go/backend/proto/iam/v1"
	"github.com/Tangyd893/ERP-Go/backend/services/iam-service/internal/app"
	apperrors "github.com/Tangyd893/ERP-Go/backend/shared/errors"
	"github.com/Tangyd893/ERP-Go/backend/shared/logger"
	"github.com/Tangyd893/ERP-Go/backend/shared/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Server 实现 iamv1.IAMServiceServer；当前仅 CheckPermission 已接入 AuthService。
type Server struct {
	iamv1.UnimplementedIAMServiceServer
	auth *app.AuthService
}

// NewServer 创建 gRPC IAM 服务实现。
func NewServer(auth *app.AuthService) *Server {
	return &Server{auth: auth}
}

// CheckPermission 委托 AuthService 校验用户权限。
func (s *Server) CheckPermission(ctx context.Context, req *iamv1.CheckPermissionRequest) (*iamv1.CheckPermissionResponse, error) {
	if req.GetUserId() == "" || req.GetTenantId() == "" || req.GetPermissionCode() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id、tenant_id、permission_code 为必填")
	}

	allowed, err := s.auth.CheckPermission(ctx, req.GetTenantId(), req.GetUserId(), req.GetPermissionCode())
	if err != nil {
		var bizErr *apperrors.BusinessError
		if errors.As(err, &bizErr) {
			switch bizErr.Code {
			case apperrors.CodeNotFound:
				return nil, status.Error(codes.NotFound, bizErr.Message)
			case apperrors.CodeUserDisabled:
				return &iamv1.CheckPermissionResponse{Allowed: false, Message: bizErr.Message}, nil
			default:
				return nil, status.Error(codes.Internal, bizErr.Message)
			}
		}
		return nil, status.Error(codes.Internal, "权限校验失败")
	}

	resp := &iamv1.CheckPermissionResponse{Allowed: allowed}
	if !allowed {
		resp.Message = "权限不足"
	}
	return resp, nil
}

func InternalTokenUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		expected := strings.TrimSpace(os.Getenv(middleware.EnvInternalToken))
		// Development may run without a token; non-dev must have one (RequireInternalTokenConfigured).
		if expected == "" {
			return handler(ctx, req)
		}
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}
		vals := md.Get(strings.ToLower(middleware.HeaderInternalToken))
		if len(vals) == 0 {
			vals = md.Get(middleware.HeaderInternalToken)
		}
		if len(vals) == 0 || vals[0] != expected {
			return nil, status.Error(codes.Unauthenticated, "invalid internal token")
		}
		return handler(ctx, req)
	}
}

// TryStart 在设置 IAM_GRPC_ADDR 时启动 gRPC 服务（默认建议 :9091）。
func TryStart(auth *app.AuthService, log logger.Logger) (stop func()) {
	addr := os.Getenv("IAM_GRPC_ADDR")
	if addr == "" {
		return func() {}
	}
	if auth == nil {
		log.Warn("IAM_GRPC_ADDR 已设置但 AuthService 未就绪，跳过 gRPC")
		return func() {}
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Errorf("IAM gRPC 监听失败 (%s): %v", addr, err)
		return func() {}
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(InternalTokenUnaryServerInterceptor()))
	iamv1.RegisterIAMServiceServer(grpcServer, NewServer(auth))

	go func() {
		if serveErr := grpcServer.Serve(lis); serveErr != nil {
			log.Errorf("IAM gRPC 服务退出: %v", serveErr)
		}
	}()

	log.Infof("IAM gRPC 已监听 %s（CheckPermission；需 %s）", addr, middleware.HeaderInternalToken)
	return func() {
		grpcServer.GracefulStop()
	}
}
