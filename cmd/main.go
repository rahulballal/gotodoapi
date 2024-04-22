package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"

	"github.com/rahulballal/gotodoapi/internal"
)

func main() {
	configPtr := internal.LoadConfig()
	mux := echo.New()
	mux.Use(middleware.Logger())

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(configPtr.LogLevel)

	internal.InitializePersistence(&log.Logger)
	todosDb := internal.NewTodosDb(&log.Logger)
	handlerMap := &internal.HandlerMap{
		Logger:      &log.Logger,
		Persistence: &todosDb,
	}
	internal.ConfigureRouting(mux, handlerMap)
	address := fmt.Sprintf(":%d", configPtr.Port)
	log.Info().Msgf("Starting server on %s", address)
	serverInitError := http.ListenAndServe(address, mux)
	log.Fatal().Err(serverInitError).Msg("Failed to start server")
}
