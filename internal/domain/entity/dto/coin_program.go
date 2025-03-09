package dto

import "github.com/google/uuid"

type CoinProgramCreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=30,plaintext"`
	Description string `json:"description" validate:"omitempty,lte=150,plaintext"`
	DayLimit    uint   `json:"day_limit" validate:"required,gte=1"`
	CardColor   string `json:"card_color" validate:"required,hexcolor"`
}

type CoinProgramCreateResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DayLimit    uint      `json:"day_limit"`
	CardColor   string    `json:"card_color"`
}

type CoinProgramReturn struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DayLimit    uint      `json:"day_limit"`
	CardColor   string    `json:"card_color"`
}

type CoinProgramUpdateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=30,plaintext"`
	Description string `json:"description" validate:"omitempty,lte=150,plaintext"`
	DayLimit    uint   `json:"day_limit" validate:"required,gte=1"`
	CardColor   string `json:"card_color" validate:"required,hexcolor"`
}

type CoinProgramUpdateResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	DayLimit    uint      `json:"day_limit"`
	CardColor   string    `json:"card_color"`
}

type CoinProgramWithRewardsReturn struct {
	CoinProgram CoinProgramParticipantReturn `json:"coin_program"`
	Rewards     []RewardReturn               `json:"rewards"`
}
