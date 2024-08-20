package repository

import (
	"errors"

	"github.com/hashicorp/go-memdb"
	"github.com/willejs/ports-service/internal/domain/entity"
)

type MemDBPortRepository struct {
	db *memdb.MemDB
}

func NewMemDBPortRepository(db *memdb.MemDB) *MemDBPortRepository {
	return &MemDBPortRepository{db: db}
}

// fetch all ports from the db
func (r *MemDBPortRepository) GetAllPorts() ([]*entity.Port, error) {
	txn := r.db.Txn(false)
	defer txn.Abort()

	it, err := txn.Get("port", "id")
	if err != nil {
		return nil, err
	}

	var ports []*entity.Port
	for obj := it.Next(); obj != nil; obj = it.Next() {
		port, ok := obj.(*entity.Port)
		if !ok {
			return nil, errors.New("failed to cast to Port")
		}
		ports = append(ports, port)
	}

	return ports, nil
}

// upsert a port into the db
func (r *MemDBPortRepository) UpsertPort(port *entity.Port) error {
	txn := r.db.Txn(true)
	defer txn.Abort()

	if err := txn.Insert("port", port); err != nil {
		return err
	}

	txn.Commit()
	return nil
}
