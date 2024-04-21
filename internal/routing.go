package internal

import "net/http"

func ConfigureRouting(mux *http.ServeMux, handlerMap *HandlerMap) {
	mux.HandleFunc("GET /todos/", handlerMap.GetTodos)
	mux.HandleFunc("POST /todos/", handlerMap.NewTodo)
	mux.HandleFunc("GET /todos/{id}", handlerMap.GetTodo)
	mux.HandleFunc("PUT /todos/{id}", handlerMap.UpdateTodo)
}
