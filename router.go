package main

import "net/http"

func Route() *http.ServeMux {
	mux := http.NewServeMux()

	// user
	mux.HandleFunc("POST /register", RegisterUserHandler)
	mux.HandleFunc("POST /login", LoginUserHandler)

	// todo
	mux.HandleFunc("GET /todos", GetTodosHandler)
	mux.HandleFunc("POST /todos", CreateTodoHandler)
	mux.HandleFunc("PUT /todos/{id}", UpdateTodoHandler)
	mux.HandleFunc("DELETE /todos/{id}", DeleteTodoHandler)

	return mux
}
