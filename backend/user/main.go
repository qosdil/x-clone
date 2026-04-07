package main

import (
	"context"
	"fmt"
	"likexuser/repository"
	"likexuser/service"
	httptransport "likexuser/transport/http"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/qosdil/like-x/backend/common/auth"
	httpauth "github.com/qosdil/like-x/backend/common/http/auth"
)

var (
	pgxPool *pgxpool.Pool
)

// main is the entry point of the user service application.
// It initializes the DB connection, constructs services and HTTP handlers,
// and starts the Fiber HTTP server.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("debug: .env file not loaded, make sure environment variables are loaded through other ways")
	}

	// Set up pgx connection pool
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB_NAME"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)
	pgxPool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(fmt.Sprintf("unable to create pgx pool: %v", err))
	}
	defer pgxPool.Close()

	// Start Fiber
	app := fiber.New()
	v1 := app.Group("/v1/users")

	// Initialize the HTTP handler with the user service and repository dependencies.
	h := httptransport.NewHandler(service.NewService(auth.NewAuth(), httpauth.NewAuth(), repository.NewPgx(pgxPool)))
	v1.Post("/authenticate", h.HandleAuthenticate)
	v1.Post("/internal/authenticate", h.HandleInternalAuthenticate)
	v1.Post("/sign-up", h.HandleSignUp)

	// Start the HTTP server and log any fatal errors.
	log.Fatal(app.Listen(":" + os.Getenv("HTTP_SERVER_PORT")))
}
