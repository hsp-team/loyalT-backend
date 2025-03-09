package coin_program

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Get(ctx context.Context, businessID uuid.UUID) (*dto.CoinProgramReturn, error) {
	coinProgram, err := s.db.CoinProgram.Query().Where(
		coinprogram.HasBusinessWith(
			business.ID(businessID),
		),
	).Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrCoinProgramNotFound
	case err != nil:
		return nil, err
	}

	return &dto.CoinProgramReturn{
		ID:          coinProgram.ID,
		Name:        coinProgram.Name,
		Description: coinProgram.Description,
		DayLimit:    coinProgram.DayLimit,
		CardColor:   coinProgram.CardColor,
	}, nil
}
