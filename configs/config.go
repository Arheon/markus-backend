package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Arheon/markus-backend/internal/shared/infrastructure/helpers"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Secret   string `env:"SECRET" yaml:"secret"`
	Env      string `env:"USER_APP_ENV" yaml:"env"`
	Database struct {
		Driver string `env:"DRIVER" yaml:"driver"`
		DSN    string `env:"DSN" yaml:"dsn"`
	} `yaml:"database"`
	Server struct {
		Port string `env:"PORT" yaml:"port"`
	} `yaml:"server"`
	Auth struct {
		JwtTokenLookup         string `env:"JWT_TOKEN_LOOKUP" yaml:"jwt_token_lookup"`
		RefreshTokenCookieName string `env:"REFRESH_TOKEN_COOKIE_NAME" yaml:"refresh_token_cookie_name"`
		IdentityKey            string `env:"IDENTITY_KEY" yaml:"identity_key"`
	} `yaml:"auth"`
	Broker struct {
		KafkaDSN     string `yaml:"kafka_dsn"`
		BrokerTopics struct {
			Server  string `env:"SERVER" yaml:"server"`
			User    string `env:"USER" yaml:"user"`
			Message string `env:"MESSAGE" yaml:"message"`
		} `yaml:"broker_topics"`
	} `yaml:"broker"`
}

const EnvTesting = "testing"

var envs = []string{"dev", "prod", EnvTesting}

func NewConfig(env string) (*Config, error) {
	if !slices.Contains(envs, env) {
		return nil, errors.New("incorrect config env")
	}

	config := &Config{
		Env: env,
	}

	file := fmt.Sprintf("configs/%s.yaml", env)

	if _, err := os.Stat(file); err != nil {
		err := cleanenv.ReadEnv(config)
		if err != nil {
			return nil, err
		}
	} else {
		err := cleanenv.ReadConfig(fmt.Sprintf("configs/%s.yaml", env), config)
		if err != nil {
			return nil, err
		}
	}

	if helpers.IsDevEnv(config.Env) {
		envHelp, _ := cleanenv.GetDescription(config, nil)
		fmt.Println(envHelp)
	}

	log.Println(env, config)

	return config, nil
}
