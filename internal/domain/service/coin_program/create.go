package coin_program

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Create(ctx context.Context, req *dto.CoinProgramCreateRequest, businessID uuid.UUID) (*dto.CoinProgramCreateResponse, error) {
	coinProgram, err := s.db.CoinProgram.Create().
		SetName(req.Name).
		SetDescription(req.Description).
		SetDayLimit(req.DayLimit).
		SetCardColor(req.CardColor).
		SetBusinessID(businessID).
		Save(ctx)
	switch {
	case ent.IsConstraintError(err):
		return nil, errorz.ErrCoinProgramAlreadyExists
	case err != nil:
		return nil, err
	}

	return &dto.CoinProgramCreateResponse{
		ID:          coinProgram.ID,
		Name:        coinProgram.Name,
		Description: coinProgram.Description,
		DayLimit:    coinProgram.DayLimit,
		CardColor:   coinProgram.CardColor,
	}, nil
}
