package reward

import (
	"context"
	"loyalit/internal/adapters/config"
	"loyalit/internal/adapters/repository/clickhouse"
	"loyalit/internal/adapters/repository/postgres/ent"
)

type statisticRepository interface {
	InsertCoinBalanceChange(ctx context.Context, change clickhouse.CoinBalanceChange) error
}

type Service struct {
	qrConfig config.QRConfig

	db                  *ent.Client
	statisticRepository statisticRepository
}

func NewService(db *ent.Client, statisticRepository statisticRepository, qrConfig config.QRConfig) *Service {
	return &Service{
		qrConfig: qrConfig,

		db:                  db,
		statisticRepository: statisticRepository,
	}
}
