package reward

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Create(ctx context.Context, req *dto.RewardCreateRequest, businessID uuid.UUID) (*dto.RewardCreateResponse, error) {
	program, err := s.db.CoinProgram.Query().
		Where(
			coinprogram.HasBusinessWith(
				business.ID(businessID),
			),
		).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorz.ErrProgramNotFound
		}

		return nil, err
	}

	reward, err := s.db.Reward.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetCost(req.Cost).
		SetImageURL(req.ImageURL).
		SetCoinProgram(program).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.RewardCreateResponse{
		ID:          reward.ID,
		Name:        reward.Name,
		Description: reward.Description,
		Cost:        reward.Cost,
		ImageURL:    reward.ImageURL,
	}, nil
}
