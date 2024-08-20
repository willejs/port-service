package memdb

import (
	"log/slog"

	"github.com/hashicorp/go-memdb"
)

// NewMemDB initializes a new MemDB instance.
func NewMemDB(logger *slog.Logger) (*memdb.MemDB, error) {
	// Define the schema for the Port entity
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"port": {
				Name: "port",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Code"},
					},
				},
			},
		},
	}

	// Create a new database instance
	logger.Info("Creating new MemDB instance", slog.String("component", "infrastructure/memdb"))
	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
