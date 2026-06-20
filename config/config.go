package config

import (
	"time"

	"github.com/nocturna-ta/golib/config"
	"github.com/nocturna-ta/golib/log"
)

type (
	MainConfig struct {
		Server     ServerConfig     `yaml:"Server"`
		API        APIConfig        `yaml:"API"`
		Database   DBConfig         `yaml:"Database"`
		Encryption EncryptionConfig `yaml:"Encryption"`
		Cors       CorsConfig       `yaml:"Cors"`
		GrpcServer GrpcConfig       `yaml:"GrpcServer"`
	}

	ServerConfig struct {
		Port         uint          `yaml:"Port" env:"SERVER_PORT"`
		WriteTimeout time.Duration `yaml:"WriteTimeout" env:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `yaml:"ReadTimeout" env:"SERVER_READ_TIMEOUT"`
	}

	GrpcConfig struct {
		Port uint `yaml:"Port" env:"GRPC_PORT"`
	}

	APIConfig struct {
		BasePath      string        `yaml:"BasePath" env:"API_BASE_PATH"`
		APITimeout    time.Duration `yaml:"APITimeout" env:"API_TIMEOUT"`
		EnableSwagger bool          `yaml:"EnableSwagger" env:"ENABLE_SWAGGER" default:"false"`
	}

	DBConfig struct {
		SlaveDSN        string `yaml:"SlaveDSN" env:"DB_SLAVE_DSN"`
		MasterDSN       string `yaml:"MasterDSN" env:"DB_MASTER_DSN"`
		RetryInterval   int    `yaml:"RetryInterval" env:"DB_RETRY_INTERVAL"`
		MaxIdleConn     int    `yaml:"MaxIdleConn" env:"DB_MAX_IDLE_CONN"`
		MaxConn         int    `yaml:"MaxConn" env:"DB_MAX_CONN"`
		ConnMaxLifetime string `yaml:"ConnMaxLifetime" env:"DB_CONN_MAX_LIFETIME"`
	}

	JWTConfig struct {
		Secret string `yaml:"Secret" env:"JWT_SECRET"`
	}

	EncryptionConfig struct {
		Key string `yaml:"Key" env:"ENCRYPTION_KEY"`
	}

	CorsConfig struct {
		AllowOrigins     string `yaml:"AllowOrigins"`
		AllowMethods     string `yaml:"AllowMethods"`
		AllowHeaders     string `yaml:"AllowHeaders"`
		AllowCredentials bool   `yaml:"AllowCredentials"`
		ExposeHeaders    string `yaml:"ExposeHeaders"`
		MaxAge           int    `yaml:"MaxAge"`
	}

	RetryConfig struct {
		MaxRetry          int             `yaml:"MaxRetry"`
		RetryInitialDelay time.Duration   `yaml:"RetryInitialDelay"`
		MaxJitter         time.Duration   `yaml:"MaxJitter"`
		HandlerTimeout    time.Duration   `yaml:"HandlerTimeout"`
		BackOffConfig     []time.Duration `yaml:"BackOffConfig"`
	}
)

func ReadConfig(cfg any, configLocation string) {
	if configLocation == "" {
		configLocation = "file://config/files/config.yaml"
	}

	if err := config.ReadConfig(cfg, configLocation, true); err != nil {
		log.WithFields(log.Fields{
			"error":           err,
			"config-location": configLocation,
		}).Fatal("Failed to read config")
	}
}
