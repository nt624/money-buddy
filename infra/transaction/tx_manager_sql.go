package transaction

import (
	"context"
	"database/sql"

	"money-buddy-backend/internal/services"
)

type sqlTx struct {
	tx *sql.Tx
}

func (t *sqlTx) Commit() error {
	return t.tx.Commit()
}

func (t *sqlTx) Rollback() error {
	return t.tx.Rollback()
}

type sqlTxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) services.TxManager {
	return &sqlTxManager{db: db}
}

func (m *sqlTxManager) Begin(ctx context.Context) (services.Tx, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlTx{tx: tx}, nil
}
