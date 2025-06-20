package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/herdiansc/go-cms/docs"
	"github.com/herdiansc/go-cms/handlers"
	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/services"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func dbOpen() error {
	var err error
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return err
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Verifying Token.....")

		svc := services.NewTokenVerifyServices()
		code, res := svc.Verify(r.Header.Get("Authorization"))
		if code != http.StatusOK {
			log.Printf("Failed to verify token\n")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Store the custom data in the request context
		ctx := context.WithValue(r.Context(), models.AuthVerifyCtxKey, res.Data)
		r = r.WithContext(ctx)

		log.Printf("Auth Res: %+v\n", res)

		next.ServeHTTP(w, r)
	})
}

// @title						Article CMS Service
// @version					1.0
// @description				A cms service built uisng golang
// @contact.email				herdiansc@gmail.com
// @license.name				MIT
// @BasePath					/
// @query.collection.format	multi
func main() {
	fmt.Println("Running server")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	if err := dbOpen(); err != nil {
		log.Fatalf("Error connecting to the database: %+v\n", err)
		return
	}
	log.Printf("DB connection established: %+v\n", DB)
	DB.AutoMigrate(&models.Auth{})
	DB.AutoMigrate(&models.Article{})
	DB.AutoMigrate(&models.ArticleTag{})
	DB.AutoMigrate(&models.Tag{})
	DB.AutoMigrate(&models.TagTrendingScore{})
	DB.AutoMigrate(&models.ArticleHistory{})

	httpServer := http.NewServeMux()

	authHandlers := handlers.NewAuthHandler(DB)
	httpServer.HandleFunc("POST /auth/register", authHandlers.Register)
	httpServer.HandleFunc("POST /auth/login", authHandlers.Login)
	httpServer.Handle("GET /auth/profile", Authenticate(http.HandlerFunc(authHandlers.GetProfile)))

	articleHandlers := handlers.NewArticleHandler(DB)
	httpServer.Handle("POST /articles", Authenticate(http.HandlerFunc(articleHandlers.Create)))
	httpServer.Handle("GET /articles", Authenticate(http.HandlerFunc(articleHandlers.List)))
	httpServer.Handle("GET /articles/{uuid}", Authenticate(http.HandlerFunc(articleHandlers.Detail)))
	httpServer.Handle("GET /articles/{uuid}/histories", Authenticate(http.HandlerFunc(articleHandlers.ListHistories)))
	httpServer.Handle("DELETE /articles/{uuid}", Authenticate(http.HandlerFunc(articleHandlers.Delete)))
	httpServer.Handle("PATCH /articles/{uuid}", Authenticate(http.HandlerFunc(articleHandlers.Patch)))

	articleHistoryHandlers := handlers.NewArticleHistoryHandler(DB)
	httpServer.Handle("GET /article-histories/{uuid}", Authenticate(http.HandlerFunc(articleHistoryHandlers.Detail)))

	tagHandlers := handlers.NewTagHandler(DB)
	httpServer.Handle("POST /tags", Authenticate(http.HandlerFunc(tagHandlers.Create)))
	httpServer.Handle("GET /tags", Authenticate(http.HandlerFunc(tagHandlers.List)))
	httpServer.Handle("GET /tags/{uuid}", Authenticate(http.HandlerFunc(tagHandlers.Detail)))

	port := os.Getenv("SERVICE_PORT")
	httpServer.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port)), //The url pointing to API definition
	))

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), httpServer)
	if err != nil {
		panic(err)
	}
}
