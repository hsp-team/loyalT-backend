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
	Create(ctx context.Context, req *dto.RewardCreateRequest, businessID uuid.UUID) (*dto.RewardCreateResponse, error)
	Delete(ctx context.Context, req *dto.RewardDeleteRequest, businessID uuid.UUID) error
	List(ctx context.Context, request *dto.RewardListRequest, businessID uuid.UUID) ([]dto.RewardReturn, error)
}

type qrService interface {
	ActivateUserReward(ctx context.Context, req *dto.RewardActivateRequest, businessID uuid.UUID) (*dto.RewardReturn, error)
}

type Handler struct {
	validator     *validator.Validator
	rewardService rewardService
	qrService     qrService
}

// NewHandler creates a new instance of reward Handler
func NewHandler(validator *validator.Validator, rewardService rewardService, qrService qrService) *Handler {
	return &Handler{
		validator:     validator,
		rewardService: rewardService,
		qrService:     qrService,
	}
}

// Setup registers reward routes
func (h *Handler) Setup(group *echo.Group) {
	group.POST("", h.Create)
	group.DELETE("/:reward_id", h.Delete)
	group.GET("", h.List)
	group.POST("/activate", h.Activate)
}
