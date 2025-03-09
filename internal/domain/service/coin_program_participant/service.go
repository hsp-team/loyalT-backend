package coin_program_participant

import (
	"loyalit/internal/adapters/repository/postgres/ent"
)

type Service struct {
	db *ent.Client
}

func NewService(db *ent.Client) *Service {
	return &Service{
		db: db,
	}
}
