package routing

import (
	"github.com/labstack/echo/v4"
	"github.com/rahulballal/gotodoapi/internal/handlers"
)

func ConfigureRouting(mux *echo.Echo, handlerMap *handlers.HandlerMap) {
	mux.GET("/todos/", handlerMap.GetTodos)
	mux.POST("/todos/", handlerMap.NewTodo)
	mux.GET("/todos/:id", handlerMap.GetTodo)
	mux.PUT("/todos/:id", handlerMap.UpdateTodo)
}
