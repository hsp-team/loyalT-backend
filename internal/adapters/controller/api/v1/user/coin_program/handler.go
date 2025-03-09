package coin_program

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks coinProgramParticipantService
type coinProgramParticipantService interface {
	UserList(ctx context.Context, req dto.CoinProgramParticipantListRequest, userID uuid.UUID) ([]dto.CoinProgramParticipantReturn, error)
	Get(ctx context.Context, coinProgramParticipantID, userID uuid.UUID) (*dto.CoinProgramParticipantReturn, error)
}

type rewardService interface {
	UserListAvailable(ctx context.Context, req *dto.CoinProgramParticipantListAvailableRequest, userID uuid.UUID) ([]dto.CoinProgramWithRewardsReturn, error)
}

type Handler struct {
	validator                     *validator.Validator
	coinProgramParticipantService coinProgramParticipantService
	rewardService                 rewardService
}

// NewHandler creates a new instance of business Handler
func NewHandler(validator *validator.Validator, coinProgramParticipantService coinProgramParticipantService, rewardService rewardService) *Handler {
	return &Handler{
		validator:                     validator,
		coinProgramParticipantService: coinProgramParticipantService,
		rewardService:                 rewardService,
	}
}

// Setup registers routes
func (h *Handler) Setup(group *echo.Group) {
	group.GET("", h.List)
	group.GET("/:coin_program_participant_id", h.Get)
	group.GET("/rewards/available", h.ListAvailable)
}
