package reward

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks rewardService
type rewardService interface {
	UserList(ctx context.Context, req *dto.RewardUserListRequest, userID uuid.UUID) ([]dto.RewardReturn, error)
	Buy(ctx context.Context, req *dto.RewardBuyRequest, userID uuid.UUID) (*dto.RewardBuyResponse, error)
}

type Handler struct {
	validator     *validator.Validator
	rewardService rewardService
}

// NewHandler creates a new instance of reward Handler
func NewHandler(validator *validator.Validator, rewardService rewardService) *Handler {
	return &Handler{
		validator:     validator,
		rewardService: rewardService,
	}
}

// Setup registers reward routes
func (h *Handler) Setup(group *echo.Group) {
	group.GET("", h.List)
	group.POST("/:reward_id/buy", h.Buy)
}
