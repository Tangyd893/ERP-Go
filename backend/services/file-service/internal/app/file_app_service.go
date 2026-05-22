package app

import (
	"context"

	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/domain"
	"github.com/Tangyd893/ERP-Go/backend/services/file-service/internal/infra/repository"
)

type FileAppService struct {
	repo *repository.FileRepository
}

func NewFileAppService(repo *repository.FileRepository) *FileAppService {
	return &FileAppService{repo: repo}
}

func (s *FileAppService) Upload(ctx context.Context, f *domain.File) error {
	return s.repo.Create(ctx, f)
}

func (s *FileAppService) Download(ctx context.Context, id string) (*domain.File, error) {
	return s.repo.FindByID(ctx, id)
}
