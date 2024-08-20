package repository

import "github.com/willejs/ports-service/internal/domain/entity"

// PortRepository defines the interface for port data operations.
type PortRepository interface {
	GetAllPorts() ([]*entity.Port, error)
	UpsertPort(port *entity.Port) error
}
