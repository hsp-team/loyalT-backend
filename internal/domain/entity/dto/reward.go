package dto

import "github.com/google/uuid"

type RewardCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50,plaintext"`
	Description string `json:"description" validate:"max=150,plaintext"`
	Cost        uint   `json:"cost" validate:"required,min=1"`
	ImageURL    string `json:"image_url" validate:"required,url"`
}

type RewardCreateResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        uint      `json:"cost"`
	ImageURL    string    `json:"image_url"`
}

type RewardDeleteRequest struct {
	ID uuid.UUID `param:"reward_id" validate:"required"`
}

type RewardListRequest struct {
	Limit  int `query:"limit" validate:"gte=0"`
	Offset int `query:"offset" validate:"gte=0"`
}

type RewardUserListRequest struct {
	ID     uuid.UUID `param:"coin_program_participant_id" validate:"required"`
	Limit  int       `query:"limit" validate:"gte=0"`
	Offset int       `query:"offset" validate:"gte=0"`
}

type RewardReturn struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Cost        uint      `json:"cost"`
	ImageURL    string    `json:"image_url"`
}

type RewardBuyRequest struct {
	CoinProgramParticipantID uuid.UUID `param:"coin_program_participant_id" validate:"required"`
	ID                       uuid.UUID `param:"reward_id" validate:"required"`
}

type RewardBuyResponse struct {
	Code string `json:"code"`
}

type RewardActivateRequest struct {
	Code string `json:"code" validate:"required,len=8"`
}
