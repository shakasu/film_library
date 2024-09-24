package initial

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Port     string
	DbDriver string
	DbSource string
}

var parameters = [...]string{"PORT", "DB_DRIVER", "DB_SOURCE"}

func Cfg() (*Config, error) {
	for _, parameter := range parameters {
		param := os.Getenv(parameter)
		if param == "" {
			return nil, errors.New(fmt.Sprintf("config [%s] in empty", parameter))
		}
	}

	return &Config{
		Port:     ":" + os.Getenv("PORT"),
		DbDriver: os.Getenv("DB_DRIVER"),
		DbSource: os.Getenv("DB_SOURCE"),
	}, nil
}
