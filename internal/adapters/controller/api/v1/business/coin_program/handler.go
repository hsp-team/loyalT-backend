package coin_program

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock_service.go -package=mocks coinProgramService
type coinProgramService interface {
	Create(ctx context.Context, req *dto.CoinProgramCreateRequest, businessID uuid.UUID) (*dto.CoinProgramCreateResponse, error)
	Get(ctx context.Context, businessID uuid.UUID) (*dto.CoinProgramReturn, error)
	Update(ctx context.Context, req *dto.CoinProgramUpdateRequest, businessID uuid.UUID) (*dto.CoinProgramUpdateResponse, error)
}

type qrService interface {
	ScanUserQR(ctx context.Context, req *dto.UserQRScanRequest, businessID uuid.UUID) (*dto.UserQRScanResponse, error)
	EnrollCoin(ctx context.Context, req *dto.UserEnrollCoinRequest, businessID uuid.UUID) error
}

type Handler struct {
	validator          *validator.Validator
	coinProgramService coinProgramService
	qrService          qrService
}

// NewHandler creates a new instance of business Handler
func NewHandler(
	validator *validator.Validator,
	coinProgramService coinProgramService,
	qrService qrService,
) *Handler {
	return &Handler{
		validator:          validator,
		coinProgramService: coinProgramService,
		qrService:          qrService,
	}
}

// Setup registers routes
func (h *Handler) Setup(group *echo.Group) {
	group.POST("", h.Create)
	group.GET("", h.Get)
	group.PUT("", h.Update)
	group.GET("/scan/:code", h.Scan)
	group.POST("/scan/enroll", h.Enroll)
}
