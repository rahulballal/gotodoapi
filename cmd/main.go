package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rahulballal/gotodoapi/internal/config"
	"github.com/rahulballal/gotodoapi/internal/handlers"
	"github.com/rahulballal/gotodoapi/internal/persistence"
	"github.com/rahulballal/gotodoapi/internal/routing"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	configPtr := config.LoadConfig()
	mux := echo.New()
	mux.Use(middleware.Logger())

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(configPtr.LogLevel)

	persistence.InitializePersistence(&log.Logger)
	todosDb := persistence.NewTodosDb(&log.Logger)
	handlerMap := &handlers.HandlerMap{
		Logger:      &log.Logger,
		Persistence: &todosDb,
	}
	routing.ConfigureRouting(mux, handlerMap)
	address := fmt.Sprintf(":%d", configPtr.Port)
	log.Info().Msgf("Starting server on %s", address)
	serverInitError := http.ListenAndServe(address, mux)
	log.Fatal().Err(serverInitError).Msg("Failed to start server")
}
