package router

import (
	"net/http"

	"github.com/MayhemI7I/URL-Shortener-Service/internal/infrastructure/http/middleware"
	"github.com/gorilla/mux"
)

type Router struct {
	mux *mux.Router
	logger interfaces.Logger
}

func NewRouter() *Router {
	router := mux.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Auth(interfaces.AuthService))

	return &Router{mux: router}
}

