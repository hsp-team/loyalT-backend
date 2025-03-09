package coin_program_participant

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogramparticipant"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) UserList(ctx context.Context, req dto.CoinProgramParticipantListRequest, userID uuid.UUID) ([]dto.CoinProgramParticipantReturn, error) {
	u, err := s.db.User.Get(ctx, userID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrUserNotFound
	case err != nil:
		return nil, err
	}

	coinProgramParticipants, err := s.db.CoinProgramParticipant.Query().Where(
		coinprogramparticipant.HasUserWith(
			user.ID(u.ID),
		),
	).Order(ent.Desc("balance")).Limit(req.Limit).Offset(req.Offset).All(ctx)
	if err != nil {
		return nil, err
	}

	resp := make([]dto.CoinProgramParticipantReturn, len(coinProgramParticipants))
	for i, coinProgramParticipant := range coinProgramParticipants {
		var coinProgram *ent.CoinProgram
		coinProgram, err = coinProgramParticipant.QueryCoinProgram().Only(ctx)
		if err != nil {
			return nil, err
		}

		resp[i] = dto.CoinProgramParticipantReturn{
			ID:          coinProgramParticipant.ID,
			Name:        coinProgram.Name,
			Description: coinProgram.Description,
			Balance:     coinProgramParticipant.Balance,
			CardColor:   coinProgram.CardColor,
		}
	}

	return resp, nil
}
