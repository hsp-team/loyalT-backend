package user

import (
	"context"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/adapters/repository/postgres/ent/user"
	"loyalit/internal/domain/common/errorz"
	"loyalit/internal/domain/entity/dto"
	"loyalit/internal/domain/utils"
)

func (s *Service) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	u, err := s.db.User.Query().Where(user.Email(req.Email)).First(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, errorz.ErrUserNotFound
	case err != nil:
		return nil, err
	}

	isMatch := utils.VerifyPassword(u.Password, req.Password)
	if !isMatch {
		return nil, errorz.ErrPasswordDoesNotMatch
	}

	refreshToken, err := utils.GenerateToken(
		u.ID,
		[]byte(s.jwtConfig.UserTokenSecret()),
		s.jwtConfig.UserTokenExpiration(),
	)
	if err != nil {
		return nil, err
	}

	return &dto.UserLoginResponse{
		Token: refreshToken,
	}, nil
}
