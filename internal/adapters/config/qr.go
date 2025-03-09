package config

import "github.com/spf13/viper"

type QRConfig interface {
	CodeLength() int
}

type qrConfig struct {
	codeLength int
}

func NewQRConfig() QRConfig {
	return qrConfig{
		codeLength: viper.GetInt("backend.qr.code-length"),
	}
}

func (c qrConfig) CodeLength() int {
	return c.codeLength
}
