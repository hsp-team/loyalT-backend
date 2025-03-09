package app

import (
	"context"
	"entgo.io/ent/dialect"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"loyalit/internal/adapters/config"
	"loyalit/internal/adapters/controller/api/v1/business"
	businesscoinprogram "loyalit/internal/adapters/controller/api/v1/business/coin_program"
	businessstatistic "loyalit/internal/adapters/controller/api/v1/business/statistic"
	usercoinprogram "loyalit/internal/adapters/controller/api/v1/user/coin_program"
	"loyalit/internal/domain/service/coin_program_participant"
	"loyalit/internal/domain/service/qr"
	"loyalit/internal/domain/service/statistic"
	"time"

	businessrewardhandler "loyalit/internal/adapters/controller/api/v1/business/coin_program/reward"
	"loyalit/internal/adapters/controller/api/v1/user"
	userrewardhandler "loyalit/internal/adapters/controller/api/v1/user/coin_program/reward"
	"loyalit/internal/adapters/controller/api/validator"
	"loyalit/internal/adapters/repository/clickhouse"
	"loyalit/internal/adapters/repository/postgres/ent"
	"loyalit/internal/domain/entity/dto"
	businessservice "loyalit/internal/domain/service/business"
	coinprogramservice "loyalit/internal/domain/service/coin_program"
	rewardservice "loyalit/internal/domain/service/reward"
	userservice "loyalit/internal/domain/service/user"
	"loyalit/pkg/closer"
	"loyalit/pkg/logger"

	_ "github.com/lib/pq"
)

type userService interface {
	Register(ctx context.Context, req *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error)
	Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.UserLoginResponse, error)
	Get(ctx context.Context, userID uuid.UUID) (*dto.UserGetResponse, error)
}

type businessService interface {
	Register(ctx context.Context, req *dto.BusinessRegisterRequest) (*dto.BusinessRegisterResponse, error)
	Login(ctx context.Context, req *dto.BusinessLoginRequest) (*dto.BusinessLoginResponse, error)
	Get(ctx context.Context, businessID uuid.UUID) (*dto.BusinessReturn, error)
	Update(ctx context.Context, req *dto.BusinessUpdateRequest, businessID uuid.UUID) (*dto.BusinessReturn, error)
}

type coinProgramService interface {
	Create(ctx context.Context, req *dto.CoinProgramCreateRequest, businessID uuid.UUID) (*dto.CoinProgramCreateResponse, error)
	Get(ctx context.Context, businessID uuid.UUID) (*dto.CoinProgramReturn, error)
	Update(ctx context.Context, req *dto.CoinProgramUpdateRequest, businessID uuid.UUID) (*dto.CoinProgramUpdateResponse, error)
}

type rewardService interface {
	Create(ctx context.Context, req *dto.RewardCreateRequest, businessID uuid.UUID) (*dto.RewardCreateResponse, error)
	Delete(ctx context.Context, req *dto.RewardDeleteRequest, businessID uuid.UUID) error
	List(ctx context.Context, request *dto.RewardListRequest, businessID uuid.UUID) ([]dto.RewardReturn, error)
	UserList(ctx context.Context, req *dto.RewardUserListRequest, userID uuid.UUID) ([]dto.RewardReturn, error)
	Buy(ctx context.Context, req *dto.RewardBuyRequest, userID uuid.UUID) (*dto.RewardBuyResponse, error)
	UserListAvailable(ctx context.Context, req *dto.CoinProgramParticipantListAvailableRequest, userID uuid.UUID) ([]dto.CoinProgramWithRewardsReturn, error)
}

type coinProgramParticipantService interface {
	UserList(ctx context.Context, req dto.CoinProgramParticipantListRequest, userID uuid.UUID) ([]dto.CoinProgramParticipantReturn, error)
	Get(ctx context.Context, coinProgramParticipantID, userID uuid.UUID) (*dto.CoinProgramParticipantReturn, error)
}

type qrService interface {
	GetUserQR(ctx context.Context, userID uuid.UUID) (*dto.QRGetResponse, error)
	ActivateUserReward(ctx context.Context, req *dto.RewardActivateRequest, businessID uuid.UUID) (*dto.RewardReturn, error)
	ScanUserQR(ctx context.Context, req *dto.UserQRScanRequest, businessID uuid.UUID) (*dto.UserQRScanResponse, error)
	EnrollCoin(ctx context.Context, req *dto.UserEnrollCoinRequest, businessID uuid.UUID) error
}

type statisticRepository interface {
	InsertCoinBalanceChange(ctx context.Context, change clickhouse.CoinBalanceChange) error
	GetUserActivityStats(ctx context.Context, userID uuid.UUID) (clickhouse.UserActivityStats, error)
	GetDailyTotalUniqueUsers(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) ([]clickhouse.DailyTotalUniqueUsers, error)
	GetDailyActiveUsers(ctx context.Context, businessID uuid.UUID, startDate, endDate time.Time) ([]clickhouse.DailyActiveUsers, error)
	GetBusinessStats(ctx context.Context, businessID uuid.UUID, periodStart, periodEnd time.Time) (clickhouse.BusinessStats, error)
	GetLoyaltyProgramStats(ctx context.Context, businessID uuid.UUID, startDate time.Time, endDate time.Time) (clickhouse.LoyaltyProgramStats, error)
	GetUserBusinessQRScansCount(
		ctx context.Context,
		businessID, userID uuid.UUID,
	) (uint64, error)
}

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

// serviceProvider is a dependency container that lazily initializes dependencies.
type serviceProvider struct {
	echo *echo.Echo

	loggerConfig     config.LoggerConfig
	pgConfig         config.PGConfig
	httpConfig       config.HTTPConfig
	jwtConfig        config.JWTConfig
	clickhouseConfig config.ClickHouseConfig
	qrConfig         config.QRConfig

	logger    *logger.Logger
	validator *validator.Validator

	db         *ent.Client
	clickhouse *clickhouse.Repository

	statisticRepository *clickhouse.StatisticRepository

	userService                   userService
	businessService               businessService
	coinProgramService            coinProgramService
	rewardService                 rewardService
	coinProgramParticipantService coinProgramParticipantService
	qrService                     qrService
	statisticService              statisticService

	userHandler                *user.Handler
	businessHandler            *business.Handler
	businessCoinProgramHandler *businesscoinprogram.Handler
	businessRewardHandler      *businessrewardhandler.Handler
	userRewardHandler          *userrewardhandler.Handler
	userCoinProgramHandler     *usercoinprogram.Handler
	businessStatsHandler       *businessstatistic.Handler
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) LoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg, err := config.NewLoggerConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get logger config: %w", err))
		}
		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		s.pgConfig = config.NewPGConfig()
	}

	return s.pgConfig
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get http config: %w", err))
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

func (s *serviceProvider) JWTConfig() config.JWTConfig {
	if s.jwtConfig == nil {
		cfg, err := config.NewJWTConfig()
		if err != nil {
			panic(fmt.Errorf("failed to get jwt config: %w", err))
		}

		s.jwtConfig = cfg
	}

	return s.jwtConfig
}

func (s *serviceProvider) ClickHouseConfig() config.ClickHouseConfig {
	if s.clickhouseConfig == nil {
		s.clickhouseConfig = config.NewClickHouseConfig()
	}

	return s.clickhouseConfig
}

func (s *serviceProvider) QRConfig() config.QRConfig {
	if s.qrConfig == nil {
		s.qrConfig = config.NewQRConfig()
	}

	return s.qrConfig
}

func (s *serviceProvider) Logger() *logger.Logger {
	if s.logger == nil {
		l, err := logger.Init(logger.Config{
			Debug:        s.LoggerConfig().Debug(),
			TimeLocation: s.LoggerConfig().TimeLocation(),
			LogToFile:    s.LoggerConfig().LogToFile(),
			LogsDir:      s.LoggerConfig().LogsDir(),
		})
		if err != nil {
			panic(fmt.Errorf("failed to init logger: %w", err))
		}

		s.logger = l
		if s.LoggerConfig().Debug() {
			s.logger.Debug("Debug mode enabled")
		}
	}

	return s.logger
}

func (s *serviceProvider) Validator() *validator.Validator {
	if s.validator == nil {
		s.validator = validator.New()
	}

	return s.validator
}

func (s *serviceProvider) DB() *ent.Client {
	if s.db == nil {
		s.logger.Debugf("Connecting to database (dsn=%s)", s.PGConfig().DSN())
		db, err := ent.Open(dialect.Postgres, s.PGConfig().DSN())
		if err != nil {
			s.Logger().Panicf("failed to open database: %v", err)
		}
		client := db

		loggerCfg := s.LoggerConfig()
		if loggerCfg.Debug() {
			client = client.Debug()
		}

		if errMigrate := client.Schema.Create(
			context.Background(),
		); errMigrate != nil {
			s.Logger().Panicf("Failed to run migrations: %v", errMigrate)
		}

		closer.Add(client.Close)
		s.db = client
	}

	return s.db
}

func (s *serviceProvider) ClickHouse() *clickhouse.Repository {
	if s.clickhouse == nil {
		ch, err := clickhouse.New(clickhouse.Config{
			Host:     s.ClickHouseConfig().Host(),
			Port:     s.ClickHouseConfig().Port(),
			Database: s.ClickHouseConfig().Database(),
			Username: s.ClickHouseConfig().Username(),
			Password: s.ClickHouseConfig().Password(),
			Debug:    s.ClickHouseConfig().Debug(),
		})

		if err != nil {
			s.Logger().Panicf("failed to open database: %v", err)
		}
		s.clickhouse = ch
	}

	return s.clickhouse
}

func (s *serviceProvider) StatisticRepository() statisticRepository {
	if s.statisticRepository == nil {
		s.statisticRepository = clickhouse.NewStatisticRepository(s.ClickHouse().GetConnection())
	}

	return s.statisticRepository
}

func (s *serviceProvider) UserService() userService {
	if s.userService == nil {
		s.userService = userservice.NewService(s.JWTConfig(), s.DB())
	}

	return s.userService
}

func (s *serviceProvider) QrService() qrService {
	if s.qrService == nil {
		s.qrService = qr.NewService(s.DB(), s.StatisticRepository(), s.QRConfig().CodeLength())
	}

	return s.qrService
}

func (s *serviceProvider) StatisticService() statisticService {
	if s.statisticService == nil {
		s.statisticService = statistic.NewService(s.DB(), s.StatisticRepository())
	}

	return s.statisticService
}

func (s *serviceProvider) UserHandler() *user.Handler {
	if s.userHandler == nil {
		s.userHandler = user.NewHandler(
			s.JWTConfig(),
			s.Validator(),
			s.UserService(),
			s.QrService(),
			s.StatisticService(),
			viper.GetBool("backend.dev-mode"),
		)
	}

	return s.userHandler
}

func (s *serviceProvider) BusinessService() businessService {
	if s.businessService == nil {
		s.businessService = businessservice.NewService(s.JWTConfig(), s.DB())
	}

	return s.businessService
}

func (s *serviceProvider) BusinessHandler() *business.Handler {
	if s.businessHandler == nil {
		s.businessHandler = business.NewHandler(
			s.JWTConfig(),
			s.Validator(),
			s.BusinessService(),
			viper.GetBool("backend.dev-mode"),
		)
	}

	return s.businessHandler
}

func (s *serviceProvider) CoinProgramService() coinProgramService {
	if s.coinProgramService == nil {
		s.coinProgramService = coinprogramservice.NewService(s.DB())
	}

	return s.coinProgramService
}

func (s *serviceProvider) BusinessCoinProgramHandler() *businesscoinprogram.Handler {
	if s.businessCoinProgramHandler == nil {
		s.businessCoinProgramHandler = businesscoinprogram.NewHandler(
			s.Validator(),
			s.CoinProgramService(),
			s.QrService(),
		)
	}

	return s.businessCoinProgramHandler
}

func (s *serviceProvider) CoinProgramParticipants() coinProgramParticipantService {
	if s.coinProgramParticipantService == nil {
		s.coinProgramParticipantService = coin_program_participant.NewService(s.DB())
	}

	return s.coinProgramParticipantService
}

func (s *serviceProvider) UserCoinProgramHandler() *usercoinprogram.Handler {
	if s.userCoinProgramHandler == nil {
		s.userCoinProgramHandler = usercoinprogram.NewHandler(s.Validator(), s.CoinProgramParticipants(), s.RewardService())
	}

	return s.userCoinProgramHandler
}

func (s *serviceProvider) RewardService() rewardService {
	if s.rewardService == nil {
		s.rewardService = rewardservice.NewService(s.DB(), s.StatisticRepository(), s.QRConfig())
	}

	return s.rewardService
}

func (s *serviceProvider) BusinessRewardHandler() *businessrewardhandler.Handler {
	if s.businessRewardHandler == nil {
		s.businessRewardHandler = businessrewardhandler.NewHandler(s.Validator(), s.RewardService(), s.QrService())
	}

	return s.businessRewardHandler
}

func (s *serviceProvider) UserRewardHandler() *userrewardhandler.Handler {
	if s.userRewardHandler == nil {
		s.userRewardHandler = userrewardhandler.NewHandler(s.Validator(), s.RewardService())
	}

	return s.userRewardHandler
}

func (s *serviceProvider) BusinessStatsHandler() *businessstatistic.Handler {
	if s.businessStatsHandler == nil {
		s.businessStatsHandler = businessstatistic.NewHandler(
			s.Validator(),
			s.StatisticService(),
		)
	}

	return s.businessStatsHandler
}
