package user

import (
	"context"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"

	"golang.org/x/crypto/bcrypt"

	entuser "loyalit/internal/adapters/repository/postgres/ent/user"
)

func (s *Service) Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	_, err := s.db.User.Query().Where(entuser.Email(req.Email)).First(ctx)
	switch {
	case err == nil:
		return nil, errorz.ErrEmailAlreadyExists
	case !ent.IsNotFound(err):
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	req.Password = string(hashedPassword)

	createdUser, err := s.db.User.Create().
		SetName(req.Name).
		SetEmail(req.Email).
		SetPassword(req.Password).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.UserRegisterResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}, nil
}
