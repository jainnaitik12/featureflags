package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"featureflags/audit-service/internal/model"
)

// AuditPostgres persists audit rows in PostgreSQL.
type AuditPostgres struct {
	pool *pgxpool.Pool
}

func NewAuditPostgres(pool *pgxpool.Pool) *AuditPostgres {
	return &AuditPostgres{pool: pool}
}

func (r *AuditPostgres) EnsureSchema(ctx context.Context) error {
	schema := `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS audit_events (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    flag_name text NOT NULL,
    action text NOT NULL,
    old_value bool NOT NULL,
    new_value bool NOT NULL,
    changed_at timestamptz NOT NULL DEFAULT now()
);`
	_, err := r.pool.Exec(ctx, schema)
	return err
}

func (r *AuditPostgres) Insert(ctx context.Context, event model.AuditEvent) error {
	_, err := r.pool.Exec(
		ctx,
		`INSERT INTO audit_events (flag_name, action, old_value, new_value)
		 VALUES ($1, $2, $3, $4)`,
		event.FlagName,
		event.Action,
		event.OldValue,
		event.NewValue,
	)
	return err
}

func (r *AuditPostgres) ListRecent(ctx context.Context, limit int) ([]model.AuditEvent, error) {
	rows, err := r.pool.Query(
		ctx,
		`SELECT id, flag_name, action, old_value, new_value, changed_at
         FROM audit_events
         ORDER BY changed_at DESC
         LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []model.AuditEvent
	for rows.Next() {
		var event model.AuditEvent
		if err := rows.Scan(
			&event.ID,
			&event.FlagName,
			&event.Action,
			&event.OldValue,
			&event.NewValue,
			&event.ChangedAt,
		); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}
