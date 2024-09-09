package repository_test

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/willejs/ports-service/internal/adapter/repository"
	"github.com/willejs/ports-service/internal/domain/entity"
	"github.com/willejs/ports-service/internal/infrastructure/memdb"
)

func TestMemDBPortRepository_GetAllPorts(t *testing.T) {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db, err := memdb.NewMemDB(logger)
	assert.NoError(t, err)
	// Create the repository
	repo := repository.NewMemDBPortRepository(db)

	// Create a new port
	port := &entity.Port{
		Name:        "Port of London",
		City:        "London",
		Country:     "United Kingdom",
		Coordinates: []float64{51.5074, 0.1278},
		Province:    "Greater London",
		Timezone:    "Europe/London",
		Unlocs:      []string{"GBLGP"},
		Code:        "GBLGP",
	}

	uErr := repo.UpsertPort(port)
	assert.NoError(t, uErr)

	// Call the GetAllPorts method
	ports, err := repo.GetAllPorts(ctx)

	// Check if there was an error
	assert.NoError(t, err)

	assert.Len(t, ports, 1)

	// Check if the returned ports are correct
	expectedPorts := []*entity.Port{
		{
			Name:        "Port of London",
			City:        "London",
			Country:     "United Kingdom",
			Coordinates: []float64{51.5074, 0.1278},
			Province:    "Greater London",
			Timezone:    "Europe/London",
			Unlocs:      []string{"GBLGP"},
			Code:        "GBLGP",
		},
	}
	assert.Equal(t, expectedPorts, ports)
}
