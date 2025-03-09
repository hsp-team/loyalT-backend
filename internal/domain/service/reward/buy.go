package reward

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/clickhouse"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogramparticipant"
	reward2 "loyalit/internal/adapters/repository/postgres/ent/reward"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/utils"
	"loyalit/internal/domain/utils/location"
	"time"
)

func (s *Service) Buy(ctx context.Context, req *dto.RewardBuyRequest, userID uuid.UUID) (*dto.RewardBuyResponse, error) {
	coinProgramParticipant, err := s.db.CoinProgramParticipant.Query().
		Where(
			coinprogramparticipant.And(
				coinprogramparticipant.ID(req.CoinProgramParticipantID),
				coinprogramparticipant.HasUserWith(
					user.ID(userID),
				),
			),
		).WithUser().WithCoinProgram().Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrCoinProgramParticipantNotFound
	case err != nil:
		return nil, err
	}

	reward, err := s.db.Reward.Get(ctx, req.ID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrRewardNotFound
	case err != nil:
		return nil, err
	}

	_, err = coinProgramParticipant.Edges.CoinProgram.QueryRewards().
		Where(
			reward2.ID(reward.ID),
		).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errorz.ErrAccessDenied
		}
		return nil, err
	}

	if coinProgramParticipant.Balance < reward.Cost {
		return nil, errorz.ErrNotEnoughCoins
	}

	tx, err := s.db.Tx(ctx)
	if err != nil {
		return nil, err
	}

	err = tx.CoinProgramParticipant.Update().
		Where(
			coinprogramparticipant.ID(coinProgramParticipant.ID),
		).
		SetBalance(coinProgramParticipant.Balance - reward.Cost).
		Exec(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	business, err := coinProgramParticipant.Edges.CoinProgram.QueryBusiness().Only(ctx)
	if err != nil {
		_ = tx.Rollback()
		if ent.IsNotFound(err) {
			return nil, errorz.ErrBusinessNotFound
		}
		return nil, err
	}

	code, err := utils.GenerateCode(s.qrConfig.CodeLength())
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	createdCoupon, err := tx.Coupon.Create().
		SetCode(code).
		SetRewardID(reward.ID).
		SetOwnerID(coinProgramParticipant.Edges.User.ID).
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = s.statisticRepository.InsertCoinBalanceChange(ctx, clickhouse.CoinBalanceChange{
		BusinessID:    business.ID,
		UserID:        coinProgramParticipant.Edges.User.ID,
		ProgramID:     coinProgramParticipant.Edges.CoinProgram.ID,
		BalanceChange: -int64(reward.Cost),
		Reason:        clickhouse.ReasonCouponBuy,
		CouponID:      &createdCoupon.ID,
		Timestamp:     time.Now().In(location.Location()),
	})
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	_ = tx.Commit()

	return &dto.RewardBuyResponse{
		Code: createdCoupon.Code,
	}, nil
}
