package internal

import (
	"github.com/rs/zerolog"
	"os"
	"strconv"
)

type Config struct {
	Port     uint64
	LogLevel zerolog.Level
}

func LoadConfig() *Config {
	Port, _ := strconv.ParseUint(os.Getenv("PORT"), 10, 32)
	LogLevelStr := os.Getenv("LOG_LEVEL")
	if Port == 0 {
		Port = 3333
	}
	return &Config{
		Port:     Port,
		LogLevel: deriveLogLevel(LogLevelStr),
	}
}

func deriveLogLevel(str string) zerolog.Level {
	level, err := zerolog.ParseLevel(str)
	if err != nil {
		return zerolog.InfoLevel
	}
	if level == zerolog.NoLevel {
		return zerolog.InfoLevel
	}
	return level
}
