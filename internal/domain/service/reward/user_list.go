package reward

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	entbusiness "loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/reward"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) UserList(ctx context.Context, req *dto.RewardUserListRequest, userID uuid.UUID) ([]dto.RewardReturn, error) {
	coinProgramParticipant, err := s.db.CoinProgramParticipant.Get(ctx, req.ID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrCoinProgramParticipantNotFound
	case err != nil:
		return nil, err
	}

	user, err := coinProgramParticipant.QueryUser().First(ctx)
	if err != nil {
		return nil, err
	}

	if user.ID != userID {
		return nil, errorz.ErrAccessDenied
	}

	coinProgram, err := coinProgramParticipant.QueryCoinProgram().Only(ctx)
	if err != nil {
		return nil, err
	}

	business, err := coinProgram.QueryBusiness().Only(ctx)
	if err != nil {
		return nil, err
	}

	rewards, err := s.db.Reward.Query().
		Where(
			reward.HasCoinProgramWith(
				coinprogram.HasBusinessWith(
					entbusiness.ID(business.ID),
				),
			),
		).
		Order(ent.Asc(reward.FieldCost)).
		Offset(req.Offset).
		Limit(req.Limit).
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
