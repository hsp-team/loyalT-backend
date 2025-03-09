package app

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"loyalit/internal/adapters/controller/api/server"
	"loyalit/pkg/closer"
	"sync"
)

// App represents the main application structure.
type App struct {
	serviceProvider *serviceProvider
	echo            *echo.Echo
}

// NewApp initializes the application and its dependencies.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, fmt.Errorf("new app: %w", err)
	}

	return a, nil
}

// Run starts the application.
func (a *App) Run() {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if a.serviceProvider.HTTPConfig().EnabledTLS() {
			a.serviceProvider.Logger().Error(a.echo.StartTLS(
				a.serviceProvider.HTTPConfig().Address(),
				a.serviceProvider.HTTPConfig().CertFile(),
				a.serviceProvider.HTTPConfig().KeyFile(),
			))
		} else {
			a.serviceProvider.Logger().Error(a.echo.Start(
				a.serviceProvider.HTTPConfig().Address(),
			))
		}
	}()

	wg.Wait()
}

// initializes application dependencies
func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return fmt.Errorf("init deps: %w", err)
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	srv := server.NewServer(
		a.serviceProvider.HTTPConfig(),
		a.serviceProvider.Logger(),
		a.serviceProvider.UserHandler(),
		a.serviceProvider.BusinessHandler(),
		a.serviceProvider.BusinessCoinProgramHandler(),
		a.serviceProvider.BusinessRewardHandler(),
		a.serviceProvider.UserCoinProgramHandler(),
		a.serviceProvider.UserRewardHandler(),
		a.serviceProvider.BusinessStatsHandler(),
	)
	a.echo = srv

	return nil
}
