package user

import (
	"context"
	"github.com/google/uuid"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
)

func (s *Service) Get(ctx context.Context, userID uuid.UUID) (*dto.UserGetResponse, error) {
	user, err := s.db.User.Get(ctx, userID)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrUserNotFound
	case err != nil:
		return nil, err
	}

	return &dto.UserGetResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
