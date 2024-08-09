package main

import (
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"

	"chat-system/database"
	"chat-system/handlers"
	"chat-system/middlewares"
)

func main() {
	database.InitRedis()
	r := mux.NewRouter()
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	authRoutes := r.PathPrefix("/").Subrouter()
	authRoutes.Use(middlewares.TokenAuthMiddleware)
	authRoutes.Handle("/send", middlewares.TokenAuthMiddleware(http.HandlerFunc(handlers.SendMessageHandler))).Methods("POST")
	authRoutes.Handle("/messages", middlewares.TokenAuthMiddleware(http.HandlerFunc(handlers.GetMessageHistoryHandler))).Methods("GET")

	// Swagger UI handler
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	fmt.Println("Listening on 8080")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
