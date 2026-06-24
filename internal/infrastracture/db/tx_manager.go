package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"task_tracker/internal/perf"
)

// in order not to create key everytime
type txKeyType struct{}

var txKey = txKeyType{}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func GetTx(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey).(*sql.Tx)
	return tx, ok
}

func (m *TxManager) WithTx(ctx context.Context, fn func(ctx context.Context) error) error {
	defer perf.Track(ctx, "tx.WithTx_total")()
	if endpoint, ok := perf.Endpoint(ctx); ok {
		stats := m.db.Stats()
		log.Printf("[PERF] endpoint=%s step=db.stats_before open=%d in_use=%d idle=%d wait_count=%d wait_duration=%s", endpoint, stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration)
	}

	beginStart := time.Now()
	tx, err := m.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if endpoint, ok := perf.Endpoint(ctx); ok {
		log.Printf("[PERF] endpoint=%s step=tx.BeginTx duration=%s", endpoint, time.Since(beginStart))
	}
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey, tx)

	fnStart := time.Now()
	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		if endpoint, ok := perf.Endpoint(ctx); ok {
			log.Printf("[PERF] endpoint=%s step=tx.fn duration=%s", endpoint, time.Since(fnStart))
		}
		return err
	}
	if endpoint, ok := perf.Endpoint(ctx); ok {
		log.Printf("[PERF] endpoint=%s step=tx.fn duration=%s", endpoint, time.Since(fnStart))
	}

	commitStart := time.Now()
	err = tx.Commit()
	if endpoint, ok := perf.Endpoint(ctx); ok {
		log.Printf("[PERF] endpoint=%s step=tx.Commit duration=%s", endpoint, time.Since(commitStart))
		stats := m.db.Stats()
		log.Printf("[PERF] endpoint=%s step=db.stats_after open=%d in_use=%d idle=%d wait_count=%d wait_duration=%s", endpoint, stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount, stats.WaitDuration)
	}
	return err
}
