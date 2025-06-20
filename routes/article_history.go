package routes

import (
	"net/http"

	"github.com/herdiansc/go-cms/handlers"
	"github.com/herdiansc/go-cms/middlewares"
	"gorm.io/gorm"
)

func ArticleHistoryRoutes(mux *http.ServeMux, DB *gorm.DB) {
	handlerFuncs := handlers.NewArticleHistoryHandler(DB)
	mux.Handle("GET /article-histories/{uuid}", middlewares.Authenticate(http.HandlerFunc(handlerFuncs.Detail)))
}
