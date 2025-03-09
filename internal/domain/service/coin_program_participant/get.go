package coin_program_participant

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Get(ctx context.Context, coinProgramParticipantID, userID uuid.UUID) (*dto.CoinProgramParticipantReturn, error) {
	coinProgramParticipant, err := s.db.CoinProgramParticipant.Get(ctx, coinProgramParticipantID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrCoinProgramParticipantNotFound
	case err != nil:
		return nil, err
	}

	user, err := coinProgramParticipant.QueryUser().Only(ctx)
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

	return &dto.CoinProgramParticipantReturn{
		ID:          coinProgramParticipant.ID,
		Name:        coinProgram.Name,
		Description: coinProgram.Description,
		Balance:     coinProgramParticipant.Balance,
		CardColor:   coinProgram.CardColor,
	}, nil
}
