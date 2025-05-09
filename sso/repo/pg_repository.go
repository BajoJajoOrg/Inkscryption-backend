package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/BajoJajoOrg/Inkscryption-backend/sso"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) sso.Repository {
	return &repository{client: client}
}

func (r *repository) Create(ctx context.Context, user sso.User) (*int, error) {
	q := `
		INSERT INTO "user" (email, password)
		VALUES ($1, $2)
		RETURNING id;
	`
	// q := `
	// 	INSERT INTO folder (name, parent_folder_id, user_id, updated_at, created_at)
	// 	VALUES ($1, $2, $3, $4, $5)
	// 	RETURNING id;
	// `

	err := r.client.QueryRow(ctx, q, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("email already exists: %w", err)
			}
			newErr := fmt.Errorf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState(),
			)
			return nil, newErr
		}
		return nil, err
	}

	return user.ID, nil
}

func (r *repository) Get(ctx context.Context, email string) (*sso.User, error) {
	q := `
		SELECT id, email, password, updated_at, created_at
		FROM "user"
		WHERE email = $1
	`

	user := sso.User{}

	print("\n", email)

	row := r.client.QueryRow(ctx, q, email)

	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.UpdatedAt, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			newErr := fmt.Errorf(
				"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
				pgErr.Message,
				pgErr.Detail,
				pgErr.Where,
				pgErr.Code,
				pgErr.SQLState(),
			)
			return nil, newErr
		}
		return nil, err
	}

	return &user, nil
}

// func (r *repository) CreateSession(ctx context.Context, s *sso.Session) (*sso.Session, error) {

// }
