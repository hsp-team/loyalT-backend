package reward

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/reward"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) List(ctx context.Context, request *dto.RewardListRequest, businessID uuid.UUID) ([]dto.RewardReturn, error) {
	rewards, err := s.db.Reward.Query().
		Where(
			reward.HasCoinProgramWith(
				coinprogram.HasBusinessWith(
					business.ID(businessID),
				),
			),
		).
		Offset(request.Offset).
		Limit(request.Limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	var response []dto.RewardReturn
	for _, r := range rewards {
		response = append(response, dto.RewardReturn{
			ID:          r.ID,
			Name:        r.Name,
			Description: r.Description,
			Cost:        r.Cost,
			ImageURL:    r.ImageURL,
		})
	}

	return response, nil
}
