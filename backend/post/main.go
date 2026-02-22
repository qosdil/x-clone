package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	ctx     context.Context
	pgxPool *pgxpool.Pool
)

// authMiddleware is a dummy authentication middleware that extracts the user ID from the "Auth-User-ID" header.
var authMiddleware = func(c fiber.Ctx) error {
	authUserID := fiber.GetReqHeader[uint](c, "Auth-User-ID")
	if authUserID == 0 || authUserID > 100_000 { // We allow only user IDs from 1 to 100,000 for testing purposes
		return c.SendStatus(http.StatusForbidden)
	}

	// Insert user if not exists (for testing purposes)
	_, err := pgxPool.Exec(c.Context(), "INSERT INTO users (id) VALUES ($1) ON CONFLICT DO NOTHING", authUserID)
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Store in request-local context
	c.Locals("authUserID", authUserID)
	return c.Next()
}

var postLikeHandler = func(c fiber.Ctx) error {
	authUserID := c.Locals("authUserID").(uint)
	postID := c.Locals("postID").(uint)

	// Insert the like record
	sql := "INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	_, err := pgxPool.Exec(c.Context(), sql, postID, authUserID)
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	// Successful
	return nil
}

var postLikeValidator = func(c fiber.Ctx) error {
	// Validate post ID format
	postID := fiber.Params[uint](c, "id")
	if postID == 0 {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Validate post existence
	var exists bool
	err := pgxPool.QueryRow(c.Context(), "SELECT EXISTS(SELECT 1 FROM posts WHERE id = $1)", postID).Scan(&exists)
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(http.StatusNotFound)
	}

	authUserID := c.Locals("authUserID").(uint)

	// Cannot like own post
	var postOwnerID uint
	err = pgxPool.QueryRow(c.Context(), "SELECT user_id FROM posts WHERE id = $1", postID).Scan(&postOwnerID)
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	if postOwnerID == authUserID {
		return c.SendStatus(http.StatusForbidden)
	}

	// Check if user already liked the post
	var alreadyLiked bool
	sql := "SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)"
	err = pgxPool.QueryRow(c.Context(), sql, postID, authUserID).Scan(&alreadyLiked)
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	if alreadyLiked {
		return c.SendStatus(http.StatusConflict)
	}

	// Store in request-local context
	c.Locals("postID", postID)
	return c.Next()
}

func init() {
	ctx = context.Background()
	godotenv.Load()

	// Set up pgx connection pool
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("X_CLONE_POSTGRES_USER"),
		os.Getenv("X_CLONE_POSTGRES_PASSWORD"),
		os.Getenv("X_CLONE_POSTGRES_HOST"),
		os.Getenv("X_CLONE_POSTGRES_PORT"),
		os.Getenv("X_CLONE_POSTGRES_DB_NAME"),
		os.Getenv("X_CLONE_POSTGRES_SSLMODE"),
	)
	var err error
	pgxPool, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		panic(fmt.Sprintf("unable to create pgx pool: %v", err))
	}
}

func main() {
	defer pgxPool.Close()

	// Start Fiber
	app := fiber.New()
	v1 := app.Group("/v1")
	v1.Post("/posts/:id/like", authMiddleware, postLikeValidator, postLikeHandler)

	log.Fatal(app.Listen(":" + os.Getenv("X_CLONE_HTTP_SERVER_PORT")))
}
