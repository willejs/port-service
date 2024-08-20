package app

import (
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

func (s *PortService) ListAllPorts() ([]*entity.Port, error) {
	return s.repo.GetAllPorts()
}

func (s *PortService) UpsertPort(port *entity.Port) error {
	return s.repo.UpsertPort(port)
}
