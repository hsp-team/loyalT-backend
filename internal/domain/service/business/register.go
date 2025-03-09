package business

import (
	"context"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Register(ctx context.Context, req *dto.BusinessRegisterRequest) (*dto.BusinessRegisterResponse, error) {
	_, err := s.db.Business.Query().Where(business.Email(req.Email)).First(ctx)
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

	createdBusiness, err := s.db.Business.Create().
		SetName(req.Name).
		SetEmail(req.Email).
		SetPassword(req.Password).
		SetDescription(req.Description).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &dto.BusinessRegisterResponse{
		ID:          createdBusiness.ID,
		Name:        createdBusiness.Name,
		Email:       createdBusiness.Email,
		Description: createdBusiness.Description,
	}, nil
}
