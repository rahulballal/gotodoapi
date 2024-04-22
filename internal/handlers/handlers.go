package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rahulballal/gotodoapi/internal/model"
	"github.com/rahulballal/gotodoapi/internal/persistence"
	"github.com/rs/zerolog"
)

type HandlerMap struct {
	Logger      *zerolog.Logger
	Persistence persistence.TodoPersistence
}

func (h *HandlerMap) GetTodos(ctx echo.Context) error {
	all, err := h.Persistence.GetAll()
	if err != nil {
		h.Logger.Err(err).Msg("Failed to get all todos")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, all)

}

func (h *HandlerMap) GetTodo(ctx echo.Context) error {
	id := ctx.Param("id")
	oneTodo, err := h.Persistence.GetByID(id)
	if err != nil {
		h.Logger.Err(err).Msg("Failed to get todo")
		return ctx.String(http.StatusInternalServerError, err.Error())
	}
	if oneTodo == nil {
		h.Logger.Info().Str("id", id).Msg("Todo not found")
		return ctx.String(http.StatusNotFound, "Not Found")
	}
	return ctx.JSON(http.StatusOK, oneTodo)
}

func (h *HandlerMap) NewTodo(ctx echo.Context) error {
	todoRequest := new(model.NewTodoRequest)
	if err := ctx.Bind(todoRequest); err != nil {
		h.Logger.Err(err).Msg("Failed to decode request")
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	todoRequest.Status = model.StatusOpen
	newTodo, dbErr := h.Persistence.Create(todoRequest)
	if dbErr != nil {
		h.Logger.Err(dbErr).Msg("Failed to create todo")
		return ctx.String(http.StatusInternalServerError, dbErr.Error())
	}
	h.Logger.Info().Str("id", newTodo.ID).Msg("Todo created")
	insertedTodo, dbErr := h.Persistence.GetByID(newTodo.ID)
	if dbErr != nil {
		h.Logger.Err(dbErr).Msg("Failed to create todo")
		return ctx.String(http.StatusInternalServerError, dbErr.Error())
	}
	return ctx.JSON(http.StatusCreated, insertedTodo)
}

func (h *HandlerMap) UpdateTodo(ctx echo.Context) error {
	id := ctx.Param("id")
	todoRequest := new(model.NewTodoRequest)
	if err := ctx.Bind(todoRequest); err != nil {
		h.Logger.Err(err).Msg("Failed to decode request")
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	todo := model.Todo{
		ID:          id,
		Title:       todoRequest.Title,
		Description: todoRequest.Description,
		Severity:    todoRequest.Severity,
		Status:      todoRequest.Status,
	}
	updatedTodo, dbErr := h.Persistence.Update(&todo)
	if dbErr != nil {
		h.Logger.Err(dbErr).Msg("Failed to update todo")
		return ctx.String(http.StatusInternalServerError, dbErr.Error())
	}
	return ctx.JSON(http.StatusOK, updatedTodo)
}
