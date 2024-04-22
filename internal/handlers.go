package internal

import (
	"encoding/json"
	"github.com/rs/zerolog"
	"net/http"
)

type HandlerMap struct {
	Logger      *zerolog.Logger
	Persistence TodoPersistence
}

func (h *HandlerMap) GetTodos(w http.ResponseWriter, r *http.Request) {
	h.Logger.Print("GetTodos Handler executed")
	all, err := h.Persistence.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if all == nil {
		_ = json.NewEncoder(w).Encode([]string{})
		return
	}
	encodingErr := json.NewEncoder(w).Encode(all)
	if encodingErr != nil {
		h.Logger.Err(encodingErr).Msg("Failed to encode response")
		http.Error(w, encodingErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerMap) GetTodo(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	oneTodo, err := h.Persistence.GetByID(id)
	if err != nil {
		h.Logger.Err(err).Msg("Failed to get todo")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if oneTodo == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	encodingErr := json.NewEncoder(w).Encode(oneTodo)
	if encodingErr != nil {
		h.Logger.Err(encodingErr).Msg("Failed to encode response")
		http.Error(w, encodingErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerMap) NewTodo(w http.ResponseWriter, r *http.Request) {
	h.Logger.Print("NewTodo Handler executed")
	// read the request body
	var todo NewTodoRequest
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		h.Logger.Err(err).Msg("Failed to decode request")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.Status = Open

	// create the todo
	newTodo, dbErr := h.Persistence.Create(&todo)
	if dbErr != nil {
		h.Logger.Err(dbErr).Msg("Failed to create todo")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	h.Logger.Info().Str("id", newTodo.ID).Msg("Todo created")
	insertedTodo, dbErr := h.Persistence.GetByID(newTodo.ID)
	if dbErr != nil {
		h.Logger.Err(dbErr).Msg("Failed to create todo")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(&insertedTodo)
	if encodeErr != nil {
		h.Logger.Err(encodeErr).Msg("Failed to encode response")
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *HandlerMap) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	h.Logger.Print("UpdateTodo Handler executed")
	id := r.PathValue("id")
	// read the request body
	var todoReq NewTodoRequest
	err := json.NewDecoder(r.Body).Decode(&todoReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// update the todo
	todo := Todo{
		ID:          id,
		Title:       todoReq.Title,
		Description: todoReq.Description,
		Severity:    todoReq.Severity,
		Status:      todoReq.Status,
	}
	updatedTodo, dbErr := h.Persistence.Update(&todo)
	if dbErr != nil {
		h.Logger.Err(dbErr).Msg("Failed to update todo")
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encodeErr := json.NewEncoder(w).Encode(updatedTodo)
	if encodeErr != nil {
		h.Logger.Err(encodeErr).Msg("Failed to encode response")
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}
}
