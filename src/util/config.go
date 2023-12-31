package util

import (
	"path"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseUrl          string        `mapstructure:"DATABASE_URL"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

func LoadConfig(dirPath string) (config Config, err error) {
	viper.SetConfigFile(path.Join(dirPath, ".env"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
