package repository

import (
	"context"
	"errors"
	user "likexuser/model"
	"testing"

	"github.com/jackc/pgx/v5"
)

type fakeRow struct {
	id  user.ID
	err error
}

func (r fakeRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != 1 {
		return errors.New("expected one destination")
	}
	p, ok := dest[0].(*user.ID)
	if !ok {
		return errors.New("invalid destination type")
	}
	*p = r.id
	return nil
}

type fakeConn struct {
	row fakeRow
}

func (c fakeConn) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.row
}

// NewDbForTest allows injecting a mock query-row connection for unit tests.
func NewDbForTest(conn queryRowConn) Repository {
	return &db{conn: conn}
}

func TestPgx_Create_Success(t *testing.T) {
	conn := fakeConn{row: fakeRow{id: 42}}
	repo := NewDbForTest(conn)

	out, err := repo.Create(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if out.ID != 42 {
		t.Fatalf("expected ID=42, got %v", out.ID)
	}
	if out.PublicID == "" {
		t.Fatal("expected non-empty PublicID")
	}
}

func TestPgx_Create_QueryError(t *testing.T) {
	conn := fakeConn{row: fakeRow{err: pgx.ErrNoRows}}
	repo := NewDbForTest(conn)

	_, err := repo.Create(context.Background(), user.CreateInput{FullName: "John Doe", Password: "secret123"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
