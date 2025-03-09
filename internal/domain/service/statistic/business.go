package statistic

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) GetBusinessStats(
	ctx context.Context,
	businessID uuid.UUID,
	req *dto.BusinessStatsRequest,
) (*dto.BusinessStatsResponse, error) {
	stats, err := s.statisticRepository.GetBusinessStats(
		ctx,
		businessID,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	return &dto.BusinessStatsResponse{
		TotalUsers:  stats.TotalUniqueUsers,
		NewUsers:    stats.NewUsersInPeriod,
		ActiveUsers: stats.NewUsersInPeriod,
	}, nil
}

func (s *Service) GetBusinessCoinProgramStats(
	ctx context.Context,
	businessID uuid.UUID,
	req *dto.BusinessCoinProgramStatsRequest,
) (*dto.BusinessCoinProgramStatsResponse, error) {
	stats, err := s.statisticRepository.GetLoyaltyProgramStats(
		ctx,
		businessID,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	return &dto.BusinessCoinProgramStatsResponse{
		TotalPointsReceived:    stats.TotalPointsReceived,
		PointsReceivedInPeriod: stats.PointsReceivedInPeriod,
		TotalCouponsPurchased:  stats.TotalCouponsPurchased,
		CouponsInPeriod:        stats.CouponsInPeriod,
	}, nil
}

func (s *Service) GetBusinessStatsDailyTotalUsers(
	ctx context.Context,
	businessID uuid.UUID,
	req *dto.BusinessStatsDailyTotalUsersRequest,
) ([]dto.BusinessStatsDailyTotalUsersResponse, error) {
	stats, err := s.statisticRepository.GetDailyTotalUniqueUsers(
		ctx,
		businessID,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	var result []dto.BusinessStatsDailyTotalUsersResponse
	for _, stat := range stats {
		result = append(result, dto.BusinessStatsDailyTotalUsersResponse{
			Date:       stat.Date,
			TotalUsers: stat.TotalUniqueUsers,
		})
	}

	return result, nil
}

func (s *Service) GetBusinessStatsDailyActiveUsers(
	ctx context.Context,
	businessID uuid.UUID,
	req *dto.BusinessStatsDailyActiveUsersRequest,
) ([]dto.BusinessStatsDailyActiveUsersResponse, error) {
	stats, err := s.statisticRepository.GetDailyActiveUsers(
		ctx,
		businessID,
		req.StartDate, req.EndDate,
	)
	if err != nil {
		return nil, err
	}

	var result []dto.BusinessStatsDailyActiveUsersResponse
	for _, stat := range stats {
		result = append(result, dto.BusinessStatsDailyActiveUsersResponse{
			Date:        stat.Date,
			ActiveUsers: stat.ActiveUsers,
		})
	}

	return result, nil
}
