package business

import (
	"context"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/business"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/utils"
)

func (s *Service) Login(ctx context.Context, req *dto.BusinessLoginRequest) (*dto.BusinessLoginResponse, error) {
	b, err := s.db.Business.Query().Where(business.Email(req.Email)).First(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrBusinessNotFound
	case err != nil:
		return nil, err
	}

	isMatch := utils.VerifyPassword(b.Password, req.Password)
	if !isMatch {
		return nil, errorz.ErrPasswordDoesNotMatch
	}

	refreshToken, err := utils.GenerateToken(
		b.ID,
		[]byte(s.jwtConfig.BusinessTokenSecret()),
		s.jwtConfig.BusinessTokenExpiration())
	if err != nil {
		return nil, err
	}

	return &dto.BusinessLoginResponse{
		Token: refreshToken,
	}, nil
}
