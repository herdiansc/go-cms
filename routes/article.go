package routes

import (
	"net/http"

	"github.com/herdiansc/go-cms/handlers"
	"github.com/herdiansc/go-cms/middlewares"
	"gorm.io/gorm"
)

func ArticleRoutes(mux *http.ServeMux, DB *gorm.DB) {
	handlerFuncs := handlers.NewArticleHandler(DB)
	mux.Handle("POST /articles", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Create)))
	mux.Handle("GET /articles", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.List)))
	mux.Handle("GET /articles/{uuid}", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Detail)))
	mux.Handle("GET /articles/{uuid}/histories", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.ListHistories)))
	mux.Handle("DELETE /articles/{uuid}", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Delete)))
	mux.Handle("PATCH /articles/{uuid}", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Patch)))
}
