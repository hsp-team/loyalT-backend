package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/adapters/config"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks userService
//go:generate mockgen -source=../../../../config/jwt.go -destination=mocks/mock_config.go -package=mocks JWTConfig
type userService interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	Get(ctx context.Context, userID uuid.UUID) (*dto.UserGetResponse, error)
}

type qrService interface {
	GetUserQR(ctx context.Context, userID uuid.UUID) (*dto.QRGetResponse, error)
}

type statisticService interface {
	GetUserStatistics(ctx context.Context, userID uuid.UUID) (*dto.UserStatisticsResponse, error)
}

type Handler struct {
	jwtConfig        config.JWTConfig
	validator        *validator.Validator
	userService      userService
	qrService        qrService
	statisticService statisticService
	devMode          bool
}

// NewHandler creates a new instance of user Handler
func NewHandler(
	jwtConfig config.JWTConfig,
	validator *validator.Validator,
	userService userService,
	qrService qrService,
	statisticService statisticService,
	devMode bool,
) *Handler {
	return &Handler{
		jwtConfig:        jwtConfig,
		validator:        validator,
		userService:      userService,
		qrService:        qrService,
		statisticService: statisticService,
		devMode:          devMode,
	}
}

// SetupAuth registers auth routes
func (h *Handler) SetupAuth(group *echo.Group) {
	group.POST("/register", h.Register)
	group.POST("/login", h.Login)
	group.POST("/logout", h.logout)
}
