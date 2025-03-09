package business

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/adapters/config"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks businessService
//go:generate mockgen -source=../../../../config/jwt.go -destination=mocks/mock_config.go -package=mocks JWTConfig
type businessService interface {
	Register(ctx context.Context, req *dto.BusinessRegisterRequest) (*dto.BusinessRegisterResponse, error)
	Login(ctx context.Context, req *dto.BusinessLoginRequest) (*dto.BusinessLoginResponse, error)
	Get(ctx context.Context, businessID uuid.UUID) (*dto.BusinessReturn, error)
	Update(ctx context.Context, req *dto.BusinessUpdateRequest, businessID uuid.UUID) (*dto.BusinessReturn, error)
}

type Handler struct {
	jwtConfig       config.JWTConfig
	validator       *validator.Validator
	businessService businessService
	devMode         bool
}

// NewHandler creates a new instance of business Handler
func NewHandler(
	jwtConfig config.JWTConfig,
	validator *validator.Validator,
	businessService businessService,
	devMode bool,
) *Handler {
	return &Handler{
		jwtConfig:       jwtConfig,
		validator:       validator,
		businessService: businessService,
		devMode:         devMode,
	}
}

// SetupAuth registers auth routes
func (h *Handler) SetupAuth(group *echo.Group) {
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/logout", h.logout)
}
