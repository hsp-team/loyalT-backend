package qr

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/coupon"
	entreward "loyalit/internal/adapters/repository/postgres/ent/reward"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) ActivateUserReward(ctx context.Context, req *dto.RewardActivateRequest, businessID uuid.UUID) (*dto.RewardReturn, error) {
	userCoupon, err := s.db.Coupon.Query().
		Where(
			coupon.And(
				coupon.Code(req.Code),
				coupon.Activated(false),
				coupon.HasRewardWith(
					entreward.HasCoinProgramWith(
						coinprogram.HasBusinessWith(
							business.ID(businessID),
						),
					),
				),
			),
		).Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrCouponNotFound
	case err != nil:
		return nil, err
	}

	reward, err := userCoupon.QueryReward().Only(ctx)
	if err != nil {
		return nil, err
	}

	err = s.db.Coupon.Update().
		Where(
			coupon.ID(userCoupon.ID),
		).SetActivated(true).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.RewardReturn{
		ID:          reward.ID,
		Name:        reward.Name,
		Description: reward.Description,
		Cost:        reward.Cost,
		ImageURL:    reward.ImageURL,
	}, nil
}
