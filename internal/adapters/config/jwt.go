package config

import (
	"github.com/spf13/viper"
	"time"
)

type JWTConfig interface {
	UserTokenSecret() string
	UserTokenExpiration() time.Duration
	BusinessTokenSecret() string
	BusinessTokenExpiration() time.Duration
}

type jwtConfig struct {
	userTokenSecret         string
	userTokenExpiration     time.Duration
	businessTokenSecret     string
	businessTokenExpiration time.Duration
}

func NewJWTConfig() (JWTConfig, error) {
	return &jwtConfig{
		userTokenSecret:         viper.GetString("backend.jwt.user-token-secret"),
		userTokenExpiration:     viper.GetDuration("backend.jwt.user-token-expiration"),
		businessTokenSecret:     viper.GetString("backend.jwt.business-token-secret"),
		businessTokenExpiration: viper.GetDuration("backend.jwt.business-token-expiration"),
	}, nil
}

func (cfg *jwtConfig) UserTokenSecret() string {
	return cfg.userTokenSecret
}

func (cfg *jwtConfig) UserTokenExpiration() time.Duration {
	return cfg.userTokenExpiration
}

func (cfg *jwtConfig) BusinessTokenSecret() string {
	return cfg.businessTokenSecret
}

func (cfg *jwtConfig) BusinessTokenExpiration() time.Duration {
	return cfg.businessTokenExpiration
}
