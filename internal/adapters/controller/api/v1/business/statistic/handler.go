package statistic

import (
	"context"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/domain/entity/dto"
)

type statisticService interface {
	GetUserStatistics(ctx context.Context, userID uuid.UUID) (*dto.UserStatisticsResponse, error)
	GetBusinessStats(
		ctx context.Context,
		businessID uuid.UUID,
		req *dto.BusinessStatsRequest,
	) (*dto.BusinessStatsResponse, error)
	GetBusinessCoinProgramStats(
		ctx context.Context,
		businessID uuid.UUID,
		req *dto.BusinessCoinProgramStatsRequest,
	) (*dto.BusinessCoinProgramStatsResponse, error)
	GetBusinessStatsDailyTotalUsers(
		ctx context.Context,
		businessID uuid.UUID,
		req *dto.BusinessStatsDailyTotalUsersRequest,
	) ([]dto.BusinessStatsDailyTotalUsersResponse, error)
	GetBusinessStatsDailyActiveUsers(
		ctx context.Context,
		businessID uuid.UUID,
		req *dto.BusinessStatsDailyActiveUsersRequest,
	) ([]dto.BusinessStatsDailyActiveUsersResponse, error)
}

type Handler struct {
	validator        *validator.Validator
	statisticService statisticService
}

// NewHandler creates a new instance of business Handler
func NewHandler(
	validator *validator.Validator,
	statisticService statisticService,
) *Handler {
	return &Handler{
		validator:        validator,
		statisticService: statisticService,
	}
}

// SetupBusinessStats registers business stats routes
func (h *Handler) SetupBusinessStats(group *echo.Group) {
	group.GET("/general", h.Business)
	group.GET("/coin_program", h.CoinProgram)
	group.GET("/daily/total_users", h.TotalUsersDaily)
	group.GET("/daily/active_users", h.ActiveUsersDaily)
}
