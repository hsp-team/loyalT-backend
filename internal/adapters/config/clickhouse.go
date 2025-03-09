package config

import "github.com/spf13/viper"

type ClickHouseConfig interface {
	Host() string
	Port() string
	Database() string
	Username() string
	Password() string
	Debug() bool
}

type clickHouseConfig struct {
	host     string
	port     string
	database string
	username string
	password string
	debug    bool
}

func NewClickHouseConfig() ClickHouseConfig {
	return &clickHouseConfig{
		host:     viper.GetString("service.clickhouse.host"),
		port:     viper.GetString("service.clickhouse.port"),
		database: viper.GetString("service.clickhouse.database"),
		username: viper.GetString("service.clickhouse.username"),
		password: viper.GetString("service.clickhouse.password"),
		debug:    viper.GetBool("settings.debug"),
	}
}

func (c *clickHouseConfig) Host() string {
	return c.host
}

func (c *clickHouseConfig) Port() string {
	return c.port
}

func (c *clickHouseConfig) Database() string {
	return c.database
}

func (c *clickHouseConfig) Username() string {
	return c.username
}

func (c *clickHouseConfig) Password() string {
	return c.password
}

func (c *clickHouseConfig) Debug() bool {
	return c.debug
}
