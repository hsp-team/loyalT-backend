package reward

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/adapters/repository/postgres/ent/coinprogram"
	"loyalit/internal/adapters/repository/postgres/ent/reward"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Delete(ctx context.Context, req *dto.RewardDeleteRequest, businessID uuid.UUID) error {
	err := s.db.Reward.DeleteOneID(req.ID).
		Where(
			reward.And(
				reward.ID(req.ID),
				reward.HasCoinProgramWith(
					coinprogram.HasBusinessWith(
						business.ID(businessID),
					),
				),
			),
		).
		Exec(ctx)
	switch {
	case ent.IsNotFound(err):
		return errorz.ErrProgramNotFound
	case err != nil:
		return err
	default:
		return nil
	}
}
