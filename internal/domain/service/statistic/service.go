package statistic

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/clickhouse"
	"loyalit/internal/adapters/repository/postgres/ent"
	"time"
)

type statisticRepository interface {
	InsertCoinBalanceChange(ctx context.Context, change clickhouse.CoinBalanceChange) error
	GetUserActivityStats(ctx context.Context, userID uuid.UUID) (clickhouse.UserActivityStats, error)
	GetDailyTotalUniqueUsers(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) ([]clickhouse.DailyTotalUniqueUsers, error)
	GetDailyActiveUsers(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) ([]clickhouse.DailyActiveUsers, error)
	GetBusinessStats(ctx context.Context, businessID uuid.UUID, periodStart, periodEnd time.Time) (clickhouse.BusinessStats, error)
	GetLoyaltyProgramStats(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) (clickhouse.LoyaltyProgramStats, error)
}

type Service struct {
	db                  *ent.Client
	statisticRepository statisticRepository
}

func NewService(db *ent.Client, statisticRepository statisticRepository) *Service {
	return &Service{
		db:                  db,
		statisticRepository: statisticRepository,
	}
}
