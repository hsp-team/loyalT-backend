package business

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Get(ctx context.Context, businessID uuid.UUID) (*dto.BusinessReturn, error) {
	business, err := s.db.Business.Get(ctx, businessID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrBusinessNotFound
	case err != nil:
		return nil, err
	}

	return &dto.BusinessReturn{
		ID:          business.ID,
		Name:        business.Name,
		Email:       business.Email,
		Description: business.Description,
	}, nil
}
