package business

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Update(ctx context.Context, req *dto.BusinessUpdateRequest, businessID uuid.UUID) (*dto.BusinessReturn, error) {
	business, err := s.db.Business.Get(ctx, businessID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrBusinessNotFound
	case err != nil:
		return nil, err
	}

	updatedBusiness, err := business.Update().
		SetName(req.Name).
		SetDescription(req.Description).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.BusinessReturn{
		ID:          updatedBusiness.ID,
		Name:        updatedBusiness.Name,
		Email:       updatedBusiness.Email,
		Description: updatedBusiness.Description,
	}, nil
}
