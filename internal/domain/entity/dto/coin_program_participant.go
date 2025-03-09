package dto

import "github.com/google/uuid"

type CoinProgramParticipantListRequest struct {
	Limit  int `query:"limit" validate:"gte=0"`
	Offset int `query:"offset" validate:"gte=0"`
}
type CoinProgramParticipantReturn struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Balance     uint      `json:"balance"`
	CardColor   string    `json:"card_color"`
}

type CoinProgramParticipantGetRequest struct {
	ID uuid.UUID `param:"coin_program_participant_id" validate:"required"`
}

type CoinProgramParticipantListAvailableRequest struct {
	Limit  int `query:"limit" validate:"gte=0"`
	Offset int `query:"offset" validate:"gte=0"`
}
