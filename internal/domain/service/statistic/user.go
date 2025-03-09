package statistic

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) GetUserStatistics(ctx context.Context, userID uuid.UUID) (*dto.UserStatisticsResponse, error) {
	userStats, err := s.statisticRepository.GetUserActivityStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserStatisticsResponse{
		UserQrScannedCount: userStats.QRScansCount,
		CouponsBought:      userStats.CouponsBought,
	}, nil
}
