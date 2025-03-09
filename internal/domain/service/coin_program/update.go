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

func (s *Service) Update(ctx context.Context, req *dto.CoinProgramUpdateRequest, businessID uuid.UUID) (*dto.CoinProgramUpdateResponse, error) {
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

	updatedCoinProgram, err := coinProgram.Update().
		SetName(req.Name).
		SetDescription(req.Description).
		SetDayLimit(req.DayLimit).
		SetCardColor(req.CardColor).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.CoinProgramUpdateResponse{
		ID:          updatedCoinProgram.ID,
		Name:        updatedCoinProgram.Name,
		Description: updatedCoinProgram.Description,
		DayLimit:    updatedCoinProgram.DayLimit,
		CardColor:   updatedCoinProgram.CardColor,
	}, nil
}
