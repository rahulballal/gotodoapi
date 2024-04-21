package internal

import (
	"os"
	"strconv"
)

type Config struct {
	Port uint64
}

func LoadConfig() *Config {
	Port, _ := strconv.ParseUint(os.Getenv("PORT"), 10, 32)
	return &Config{
		Port: Port,
	}
}
