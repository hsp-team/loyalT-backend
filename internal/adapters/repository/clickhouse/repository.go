package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Config struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	Debug    bool
}

type Repository struct {
	conn driver.Conn
}

func New(cfg Config) (*Repository, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.Database,
			Username: cfg.Username,
			Password: cfg.Password,
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Debug: cfg.Debug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to clickhouse: %w", err)
	}

	if err := createTables(conn); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &Repository{conn: conn}, nil
}

func (r *Repository) GetConnection() driver.Conn {
	return r.conn
}

func createTables(conn driver.Conn) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS coin_balance_changes (
			business_id UUID,
			user_id UUID,
			program_id UUID,
			balance_change Int64,
			reason String,
			coupon_id Nullable(UUID),
			timestamp DateTime64(3),
			INDEX idx_business_id business_id TYPE bloom_filter GRANULARITY 1,
			INDEX idx_user_id user_id TYPE bloom_filter GRANULARITY 1,
			INDEX idx_program_id program_id TYPE bloom_filter GRANULARITY 1,
			INDEX idx_timestamp timestamp TYPE minmax GRANULARITY 1
		) ENGINE = MergeTree()
		ORDER BY (timestamp, business_id, user_id, program_id)
		PARTITION BY toYYYYMM(timestamp)`,
	}

	for _, query := range queries {
		if err := conn.Exec(context.Background(), query); err != nil {
			return fmt.Errorf("failed to execute query: %w", err)
		}
	}

	return nil
}
