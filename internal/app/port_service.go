package app

import (
	"context"
	"github.com/willejs/ports-service/internal/domain/entity"
	"github.com/willejs/ports-service/internal/domain/repository"
)

// PortService handles the business logic for ports.
type PortService struct {
	repo repository.PortRepository
}

func NewPortService(repo repository.PortRepository) *PortService {
	return &PortService{repo: repo}
}

func (s *PortService) ListAllPorts(ctx context.Context) ([]*entity.Port, error) {
	return s.repo.GetAllPorts(ctx)
}

func (s *PortService) UpsertPort(port *entity.Port) error {
	return s.repo.UpsertPort(port)
}
