package repository

import (
	"context"
	"fmt"
	user "likexuser/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"golang.org/x/crypto/bcrypt"
)

type queryRowConn interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type db struct {
	conn queryRowConn
}

// Create inserts a new user record and returns the created user output.
func (r *db) Create(ctx context.Context, input user.CreateInput) (output user.CreateOutput, err error) {
	sql := "INSERT INTO users (public_id, full_name, password_hash) VALUES ($1, $2, $3) RETURNING id"
	var id user.ID
	publicID, err := gonanoid.New()
	if err != nil {
		return user.CreateOutput{}, fmt.Errorf("failed to generate Nano ID: %v", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return user.CreateOutput{}, fmt.Errorf("failed to hash password: %v", err)
	}

	err = r.conn.QueryRow(ctx, sql, publicID, input.FullName, hash).Scan(&id)
	if err != nil {
		return user.CreateOutput{}, fmt.Errorf("failed to create user: %v", err)
	}

	return user.CreateOutput{ID: id, PublicID: user.PublicID(publicID)}, nil
}

// NewPgx creates a Repository implementation backed by a pgx connection pool.
func NewPgx(conn *pgxpool.Pool) Repository {
	return &db{conn: conn}
}
