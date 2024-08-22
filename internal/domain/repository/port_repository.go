package repository

import (
	"context"

	"github.com/willejs/ports-service/internal/domain/entity"
)

// PortRepository defines the interface for port data operations.
type PortRepository interface {
	GetAllPorts(ctx context.Context) ([]*entity.Port, error)
	UpsertPort(port *entity.Port) error
}
