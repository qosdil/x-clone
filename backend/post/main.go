package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/qosdil/x-clone/backend/common/http/auth"
)

var (
	pgxPool *pgxpool.Pool
)

var postLikeHandler = func(c fiber.Ctx) error {
	authUserID := c.Locals("auth_user_id").(uint)
	postID := c.Locals("postID").(uint)
	sql := "INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	if os.Getenv("ASYNC_POST_LIKE") == "true" {
		go func() {
			_, err := pgxPool.Exec(context.Background(), sql, postID, authUserID)
			if err != nil {
				log.Printf("async database error: %v", err)
				return
			}
			log.Printf("user %d liked post %d (async)", authUserID, postID)
		}()
		return c.SendStatus(http.StatusAccepted)
	}

	// Synchronous insert
	_, err := pgxPool.Exec(c.Context(), sql, postID, authUserID)
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	log.Printf("user %d liked post %d", authUserID, postID)
	return c.SendStatus(http.StatusOK)
}

var postLikeValidator = func(c fiber.Ctx) error {
	// Validate post ID format
	publicPostID := c.Params("public_id", "")
	if publicPostID == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	// Validate post existence
	var postID, postOwnerID uint
	sql := "SELECT id, user_id FROM posts WHERE posts.public_id = $1"
	err := pgxPool.QueryRow(c.Context(), sql, publicPostID).Scan(&postID, &postOwnerID)
	if err == pgx.ErrNoRows {
		return c.SendStatus(http.StatusNotFound)
	}
	if err != nil {
		log.Printf("database error: %v", err)
		return c.SendStatus(http.StatusInternalServerError)
	}

	authUserID := c.Locals("auth_user_id").(uint)

	// Cannot like own post
	if postOwnerID == authUserID {
		return c.SendStatus(http.StatusForbidden)
	}

	// Check if user already liked the post
	var alreadyLiked bool
	sql = "SELECT EXISTS(SELECT 1 FROM post_likes WHERE post_id = $1 AND user_id = $2)"
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
	pgxPool, err = pgxpool.New(context.Background(), dbURL)
	if err != nil {
		panic(fmt.Sprintf("unable to create pgx pool: %v", err))
	}
}

func main() {
	defer pgxPool.Close()

	// Start Fiber
	app := fiber.New()
	if os.Getenv("DEBUG") == "true" {
		app.Use(pprof.New())
	}
	v1 := app.Group("/v1")
	v1.Post("/posts/:public_id/like", auth.AuthMiddleware, postLikeValidator, postLikeHandler)

	log.Fatal(app.Listen(":" + os.Getenv("X_CLONE_HTTP_SERVER_PORT")))
}
