package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

type HandlerMap struct {
	Logger      *log.Logger
	Persistence TodoPersistence
}

func (h *HandlerMap) GetTodos(w http.ResponseWriter, r *http.Request) {
	h.Logger.Print("GetTodos Handler executed")
	all, error := h.Persistence.GetAll()
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(all)
}

func (h *HandlerMap) GetTodo(w http.ResponseWriter, r *http.Request) {
	h.Logger.Print("GetTodo Handler executed")
	id := r.PathValue("id")
	oneTodo, error := h.Persistence.GetByID(id)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(oneTodo)
}

func (h *HandlerMap) NewTodo(w http.ResponseWriter, r *http.Request) {
	h.Logger.Print("NewTodo Handler executed")
	// read the request body
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	todo.Status = Open

	// create the todo
	newTodo, error := h.Persistence.Create(&todo)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	insertedTodo, error := h.Persistence.GetByID(newTodo.ID)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&insertedTodo)
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
	updatedTodo, error := h.Persistence.Update(&todo)
	if error != nil {
		http.Error(w, error.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTodo)
}
