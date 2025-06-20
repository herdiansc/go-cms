package routes

import (
	"fmt"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

func LoadRoutes(DB *gorm.DB) http.Handler {
	httpServer := http.NewServeMux()

	AuthRoutes(httpServer, DB)
	ArticleRoutes(httpServer, DB)
	ArticleHistoryRoutes(httpServer, DB)
	TagRoutes(httpServer, DB)

	port := os.Getenv("SERVICE_PORT")
	httpServer.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port)), //The url pointing to API definition
	))

	return httpServer
}
