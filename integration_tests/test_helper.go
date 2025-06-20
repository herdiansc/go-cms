package integrationtests

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// _ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/herdiansc/go-cms/config"
	"github.com/herdiansc/go-cms/models"
	"github.com/herdiansc/go-cms/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	Port string
	// DbInstance *sqlx.DB
	DB        *gorm.DB
	Container testcontainers.Container
}

func SetupTestDatabase() TestDatabase {

	// setup db container
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, db, port, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	cancel()

	db.AutoMigrate(&models.Auth{})
	db.AutoMigrate(&models.Article{})
	db.AutoMigrate(&models.ArticleTag{})
	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.TagTrendingScore{})
	db.AutoMigrate(&models.ArticleHistory{})

	return TestDatabase{
		Port:      port,
		DB:        db,
		Container: container,
	}
}

func (tdb *TestDatabase) TearDown() {
	// tdb.DbInstance.Close()
	// remove test container
	_ = tdb.Container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *gorm.DB, string, error) {
	config.LoadEnv("../.env.integration.test")
	var env = map[string]string{
		"POSTGRES_PASSWORD": os.Getenv("DB_PASSWORD"),
		"POSTGRES_USER":     os.Getenv("DB_USER"),
		"POSTGRES_DB":       os.Getenv("DB_NAME"),
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), p.Port())

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return container, nil, p.Port(), fmt.Errorf("failed to establish database connection: %v", err)
	}

	// dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	// db, err := pgxpool.Connect(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName))
	// if err != nil {
	// 	return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	// }

	return container, db, p.Port(), nil
}

// func createContainer(ctx context.Context) (testcontainers.Container, *gorm.DB, error) {
// 	config.LoadEnv("../.env.integration.test")

// 	pgContainer, err := postgres.Run(ctx,
// 		"postgres:alpine",
// 		// postgres.WithInitScripts(migrationFilesPath...),
// 		postgres.WithDatabase(os.Getenv("DB_NAME")),
// 		postgres.WithUsername(os.Getenv("DB_USER")),
// 		postgres.WithPassword(os.Getenv("DB_PASSWORD")),
// 		testcontainers.WithWaitStrategy(
// 			wait.ForLog("database system is ready to accept connections").
// 				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
// 	)
// 	fmt.Printf("EEEEEEEEEEEEEE: %+v\n", err)
// 	if err != nil {
// 		return pgContainer, nil, err
// 	}

// 	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")

// 	if strings.Contains(strings.ToLower(connStr), "localhost") {
// 		connStr = strings.Replace(connStr, "localhost", os.Getenv("DB_HOST"), 1)
// 	}

// 	if err != nil {
// 		return pgContainer, nil, err
// 	}

// 	fmt.Printf("XXXXXX: %+v\n", err)

// 	if err != nil {
// 		return pgContainer, nil, fmt.Errorf("failed to establish database connection: %v", err)
// 	}

// 	return pgContainer, nil, nil
// }

func setupServer(dbPort string) http.Handler {
	config.LoadEnv("../.env.integration.test")
	DB := config.SetupDB(dbPort)
	return routes.LoadRoutes(DB)
}

func StartServer(dbPort string) {
	httpServer := setupServer(dbPort)

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("SERVICE_PORT")), httpServer)
	if err != nil {
		panic(err)
	}
}
