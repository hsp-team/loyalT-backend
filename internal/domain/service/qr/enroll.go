package qr

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/clickhouse"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogramparticipant"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/utils/location"
	"time"
)

func (s *Service) EnrollCoin(ctx context.Context, req *dto.UserEnrollCoinRequest, businessID uuid.UUID) error {
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return err
	}

	participant, err := tx.CoinProgramParticipant.Query().Where(
		coinprogramparticipant.And(
			coinprogramparticipant.HasUserWith(
				user.QrData(req.Code),
			),
			coinprogramparticipant.HasCoinProgramWith(
				coinprogram.HasBusinessWith(
					business.ID(businessID),
				),
			),
		),
	).WithCoinProgram().WithUser().Only(ctx)
	if err != nil {
		_ = tx.Rollback()
		switch {
		case ent.IsNotFound(err):
			return errorz.ErrCoinProgramParticipantNotFound
		}
		return err
	}

	err = participant.Update().AddBalance(1).Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return errorz.ErrCoinProgramParticipantNotFound
	}

	err = tx.User.Update().Where(
		user.QrData(req.Code),
	).SetQrData("").Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		switch {
		case ent.IsNotFound(err):
			return errorz.ErrUserByQrNotFound
		}
		return err
	}

	scanCount, err := s.statisticRepository.GetUserBusinessQRScansCount(
		ctx,
		businessID,
		participant.Edges.User.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if scanCount >= uint64(participant.Edges.CoinProgram.DayLimit) {
		_ = tx.Rollback()
		return errorz.ErrUserScanLimitReached
	}

	err = s.statisticRepository.InsertCoinBalanceChange(ctx, clickhouse.CoinBalanceChange{
		BusinessID:    businessID,
		UserID:        participant.Edges.User.ID,
		ProgramID:     participant.Edges.CoinProgram.ID,
		BalanceChange: 1,
		Reason:        clickhouse.ReasonQRScan,
		CouponID:      nil,
		Timestamp:     time.Now().In(location.Location()),
	})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
