package qr

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogramparticipant"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) ScanUserQR(ctx context.Context, req *dto.UserQRScanRequest, businessID uuid.UUID) (*dto.UserQRScanResponse, error) {
	u, err := s.db.User.Query().Where(
		user.QrData(req.Code),
	).Only(ctx)
	if err != nil {
		return nil, errorz.ErrUserByQrNotFound
	}

	var participant *ent.CoinProgramParticipant
	participant, err = s.db.CoinProgramParticipant.Query().Where(
		coinprogramparticipant.And(
			coinprogramparticipant.HasUserWith(
				user.ID(u.ID),
			),
			coinprogramparticipant.HasCoinProgramWith(
				coinprogram.HasBusinessWith(
					business.ID(businessID),
				),
			),
		),
	).Only(ctx)
	if err != nil {
		coinProgram, err := s.db.CoinProgram.Query().Where(
			coinprogram.HasBusinessWith(
				business.ID(businessID),
			),
		).Only(ctx)
		if err != nil {
			return nil, errorz.ErrCoinProgramNotFound
		}

		participant, err = s.db.CoinProgramParticipant.Create().
			SetUser(u).
			SetCoinProgram(coinProgram).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &dto.UserQRScanResponse{
		Username: u.Name,
		Balance:  participant.Balance,
	}, nil
}
