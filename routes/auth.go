package routes

import (
	"net/http"

	"github.com/herdiansc/go-cms/handlers"
	"github.com/herdiansc/go-cms/middlewares"
	"gorm.io/gorm"
)

func AuthRoutes(mux *http.ServeMux, DB *gorm.DB) {
	handlerFuncs := handlers.NewAuthHandler(DB)
	mux.HandleFunc("POST /auth/register", handlerFuncs.Register)
	mux.HandleFunc("POST /auth/login", handlerFuncs.Login)
	mux.Handle("GET /auth/profile", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.GetProfile)))
}
