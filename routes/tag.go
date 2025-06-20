package routes

import (
	"net/http"

	"github.com/herdiansc/go-cms/handlers"
	"github.com/herdiansc/go-cms/middlewares"
	"gorm.io/gorm"
)

func TagRoutes(mux *http.ServeMux, DB *gorm.DB) {
	handlerFuncs := handlers.NewTagHandler(DB)
	mux.Handle("GET /tags", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.List)))
	mux.Handle("GET /tags/{id}", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Detail)))
	mux.Handle("POST /tags", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Create)))
}
