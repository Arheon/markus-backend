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
	Env      string `env:"USER_APP_ENV"`
	Database struct {
		Driver string `env:"DRIVER"`
		DSN    string `env:"DSN"`
	} `env-prefix:"USER_APP_DB_"`
	Server struct {
		Port string `env:"PORT"`
	}
	Auth struct {
		JwtTokenLookup         string `env:"JWT_TOKEN_LOOKUP"`
		RefreshTokenCookieName string `env:"REFRESH_TOKEN_COOKIE_NAME"`
		IdentityKey            string `env:"IDENTITY_KEY"`
	} `env-prefix:"USER_APP_AUTH_"`
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

	log.Println(config)

	return config, nil
}
