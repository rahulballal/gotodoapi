package internal

import (
	"github.com/labstack/echo/v4"
)

func ConfigureRouting(mux *echo.Echo, handlerMap *HandlerMap) {
	mux.GET("/todos/", handlerMap.GetTodos)
	mux.POST("/todos/", handlerMap.NewTodo)
	mux.GET("/todos/:id", handlerMap.GetTodo)
	mux.PUT("/todos/:id", handlerMap.UpdateTodo)
}
