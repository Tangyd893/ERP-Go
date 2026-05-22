package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/tenant-service/internal/infra/repository"
)

// TenantAppService 租户应用服务
type TenantAppService struct {
	tenantRepo *repository.TenantRepository
	orgRepo    *repository.OrgRepository
}

// NewTenantAppService 创建租户应用服务
func NewTenantAppService(tenantRepo *repository.TenantRepository, orgRepo *repository.OrgRepository) *TenantAppService {
	return &TenantAppService{
		tenantRepo: tenantRepo,
		orgRepo:    orgRepo,
	}
}

func (s *TenantAppService) ListTenants(ctx context.Context, offset, limit int) ([]*domain.Tenant, int64, error) {
	return s.tenantRepo.List(ctx, offset, limit)
}

func (s *TenantAppService) GetTenant(ctx context.Context, id string) (*domain.Tenant, error) {
	return s.tenantRepo.FindByID(ctx, id)
}

func (s *TenantAppService) CreateTenant(ctx context.Context, tenant *domain.Tenant) error {
	return s.tenantRepo.Create(ctx, tenant)
}

func (s *TenantAppService) UpdateTenant(ctx context.Context, tenant *domain.Tenant) error {
	return s.tenantRepo.Update(ctx, tenant)
}

func (s *TenantAppService) ListOrganizations(ctx context.Context, tenantID string) ([]*domain.Organization, error) {
	return s.orgRepo.ListOrgsByTenant(ctx, tenantID)
}

func (s *TenantAppService) CreateOrganization(ctx context.Context, org *domain.Organization) error {
	return s.orgRepo.CreateOrg(ctx, org)
}

func (s *TenantAppService) ListDepartments(ctx context.Context, orgID string) ([]*domain.Department, error) {
	return s.orgRepo.ListDeptsByOrg(ctx, orgID)
}

func (s *TenantAppService) CreateDepartment(ctx context.Context, dept *domain.Department) error {
	return s.orgRepo.CreateDept(ctx, dept)
}

func (s *TenantAppService) ListPositions(ctx context.Context, deptID string) ([]*domain.Position, error) {
	return s.orgRepo.ListPositionsByDept(ctx, deptID)
}

func (s *TenantAppService) CreatePosition(ctx context.Context, pos *domain.Position) error {
	return s.orgRepo.CreatePosition(ctx, pos)
}
