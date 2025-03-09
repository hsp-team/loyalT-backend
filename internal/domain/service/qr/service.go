package qr

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/clickhouse"
	"loyalit/internal/adapters/repository/postgres/ent"
)

type statisticRepository interface {
	InsertCoinBalanceChange(ctx context.Context, change clickhouse.CoinBalanceChange) error
	GetUserBusinessQRScansCount(ctx context.Context, businessID, userID uuid.UUID) (uint64, error)
}

type Service struct {
	db                  *ent.Client
	statisticRepository statisticRepository
	codeLength          int
}

func NewService(db *ent.Client, statisticRepository statisticRepository, codeLength int) *Service {
	return &Service{
		db:                  db,
		statisticRepository: statisticRepository,
		codeLength:          codeLength,
	}
}
