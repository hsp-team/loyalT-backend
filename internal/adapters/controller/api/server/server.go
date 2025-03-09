package server

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"golang.org/x/time/rate"
	"loyalit/internal/adapters/config"
	"loyalit/internal/adapters/controller/api/v1/business"
	businesscoinprogram "loyalit/internal/adapters/controller/api/v1/business/coin_program"
	businessreward "loyalit/internal/adapters/controller/api/v1/business/coin_program/reward"
	businessstatistic "loyalit/internal/adapters/controller/api/v1/business/statistic"
	"loyalit/internal/adapters/controller/api/v1/user"
	usercoinprogram "loyalit/internal/adapters/controller/api/v1/user/coin_program"
	userreward "loyalit/internal/adapters/controller/api/v1/user/coin_program/reward"
	"loyalit/internal/domain/entity/dto"
	"loyalit/pkg/logger"
)

// NewServer initializes and returns a new server.
func NewServer(
	httpConfig config.HTTPConfig,
	logger *logger.Logger,
	userHandler *user.Handler,
	businessHandler *business.Handler,
	businessCoinProgramHandler *businesscoinprogram.Handler,
	businessRewardHandler *businessreward.Handler,
	userCoinProgramHandler *usercoinprogram.Handler,
	userRewardHandler *userreward.Handler,
	businessStatsHandler *businessstatistic.Handler,
) *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := echo.ErrBadRequest.Code
		var he *echo.HTTPError
		if errors.As(err, &he) {
			code = he.Code
		}

		if !c.Response().Committed {
			_ = c.JSON(code, dto.HTTPStatus{
				Code:    code,
				Message: err.Error(),
			})
		}
	}

	e.Use(middleware.Secure())
	//e.Use(middleware.CSRF())
	e.Use(middleware.Recover())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))
	e.Use(middleware.CORSWithConfig(httpConfig.GetCORSConfig()))

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		HandleError: true,
		LogError:    true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Infow("request completed",
					"ip", v.RemoteIP,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
				)
			} else {
				logger.Errorw("request failed",
					"ip", v.RemoteIP,
					"method", v.Method,
					"uri", v.URI,
					"status", v.Status,
					"error", v.Error.Error(),
				)
			}
			return nil
		},
	}))

	addRoutes(
		e,
		logger,
		userHandler,
		businessHandler,
		businessCoinProgramHandler,
		businessRewardHandler,
		userCoinProgramHandler,
		userRewardHandler,
		businessStatsHandler,
	)

	return e
}

// registers application routes to the provided router.
func addRoutes(
	e *echo.Echo,
	logger *logger.Logger,
	userHandler *user.Handler,
	businessHandler *business.Handler,
	businessCoinProgramHandler *businesscoinprogram.Handler,
	businessRewardHandler *businessreward.Handler,
	userCoinProgramHandler *usercoinprogram.Handler,
	userRewardHandler *userreward.Handler,
	businessStatsHandler *businessstatistic.Handler,
) {
	logger.Debug("Setting up routes")

	api := e.Group("/api")
	api.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, "pong")
	})
	v1 := api.Group("/v1")
	v1.GET("/docs/*", echoSwagger.WrapHandler)

	u := v1.Group("/user")
	b := v1.Group("/business")

	userHandler.SetupAuth(u.Group("/auth"))
	businessHandler.SetupAuth(b.Group("/auth"))
	u.Use(userHandler.Auth)
	u.GET("/me", userHandler.Me)
	u.GET("/qr", userHandler.GetQR)
	b.Use(businessHandler.Auth)
	b.GET("/me", businessHandler.Me)
	b.PUT("", businessHandler.Update)

	bcp := b.Group("/coin_program")
	businessCoinProgramHandler.Setup(bcp)
	businessRewardHandler.Setup(bcp.Group("/rewards"))

	ucp := u.Group("/coin_programs")
	userCoinProgramHandler.Setup(ucp)
	userRewardHandler.Setup(ucp.Group("/:coin_program_participant_id/rewards"))

	businessStatsHandler.SetupBusinessStats(b.Group("/stats"))

	routes, err := json.MarshalIndent(e.Routes(), "", "  ")
	if err == nil {
		logger.Debugf("Routes: %s", routes)
	}
}
