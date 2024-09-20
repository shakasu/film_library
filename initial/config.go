package initial

import (
	"os"
)

type Config struct {
	Port     string
	DbDriver string
	DbSource string
}

func Cfg() Config {
	return Config{
		Port:     ":" + os.Getenv("PORT"),
		DbDriver: os.Getenv("DB_DRIVER"),
		DbSource: os.Getenv("DB_SOURCE"),
	}
}
