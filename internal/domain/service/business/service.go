package business

import (
	"loyalit/internal/adapters/config"
	"loyalit/internal/adapters/repository/postgres/ent"
)

//go:generate mockgen -source=../../../adapters/config/jwt.go -destination=mocks/mock_config.go -package=mocks JWTConfig
type Service struct {
	jwtConfig config.JWTConfig

	db *ent.Client
}

func NewService(jwtConfig config.JWTConfig, db *ent.Client) *Service {
	return &Service{
		jwtConfig: jwtConfig,

		db: db,
	}
}
