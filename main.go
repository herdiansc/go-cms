package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/herdiansc/go-cms/config"
	_ "github.com/herdiansc/go-cms/docs"
	"github.com/herdiansc/go-cms/routes"
)

func setupServer() http.Handler {
	config.LoadEnv(".env")
	DB := config.SetupDB("")
	return routes.LoadRoutes(DB)
}

func main() {
	fmt.Println("Server Running")
	httpServer := setupServer()

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")), httpServer)
	if err != nil {
		panic(err)
	}
}
