package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/db/postgresql"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// func NewClient(ctx context.Context, maxAttempts int, sc config.) {
// 	dsn := fmt.Sprintf
// }

type repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) canvas.Repository {
	return &repository{client: client}
}

func (r *repository) Update(ctx context.Context, id int, url string) error {

	var currentURL sql.NullString
	q := `
		SELECT url
		FROM canvas
		WHERE id = $1
	`
	err := r.client.QueryRow(ctx, q, id).Scan(&currentURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("canvas with such id was not found")
		}
		return fmt.Errorf("sql error")
	}

	if currentURL.Valid && currentURL.String != "" {
		return nil
	}

	q = `
		UPDATE canvas
		SET url = $1
		WHERE id = $2
	`

	_, err = r.client.Exec(ctx, q, url, id)
	if err != nil {
		return fmt.Errorf("error while updating canvas url")
	}
	return nil
}

func (r *repository) UpdateText(ctx context.Context, id int, text string) error {

	q := `
		UPDATE canvas
		SET text = $1
		WHERE id = $2
	`

	_, err := r.client.Exec(ctx, q, text, id)
	if err != nil {
		return fmt.Errorf("error while updating canvas url")
	}

	return nil
}

// Create implements canvas.CanvasStorage.
func (r *repository) Create(ctx context.Context, canvas canvas.CanvasBase) (*int, error) {
	q := `
		INSERT INTO canvas (name, updated_at, folder_id, user_id, created_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id;
		`
	err := r.client.QueryRow(ctx, q, canvas.Name, canvas.UpdatedAt, canvas.FolderId, canvas.UserId, canvas.CreatedAt).Scan(&canvas.CanvasID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
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
	return canvas.CanvasID, nil
}

func (r *repository) GetAll(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) ([]canvas.CanvasBase, error) {

	// сделать нормальные ошибки
	// TODO: подумать как сделать нормальные args
	// fmt.Println("Ваш user_id: ", id)

	query, args, err := postgresql.BuildQuery(filterOptions, folder_id, user_id)
	if err != nil {
		return nil, err
	}

	pgxArgs := make([]any, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		// case time.Time:
		// 	pgxArgs = append(pgxArgs, v.Format(time.RFC3339))
		default:
			pgxArgs = append(pgxArgs, v)
		}
	}

	rows, err := r.client.Query(ctx, query, pgxArgs...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
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

	canvases := make([]canvas.CanvasBase, 0)

	var nullUrl sql.NullString
	var nullText sql.NullString

	for rows.Next() {
		var canvas canvas.CanvasBase

		err = rows.Scan(&canvas.CanvasID, &canvas.Name, &nullUrl, &canvas.UpdatedAt, &nullText, &canvas.FolderId, &canvas.UserId, &canvas.CreatedAt)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.Is(err, pgErr) {
				pgErr = err.(*pgconn.PgError)
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

		if nullUrl.Valid {
			canvas.Url = nullUrl.String
		} else {
			canvas.Url = ""
		}

		if nullText.Valid {
			canvas.Text = nullText.String
		} else {
			canvas.Text = ""
		}

		canvases = append(canvases, canvas)
	}

	if err = rows.Err(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			// newErr := fmt.Errorf(
			// 	"SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s",
			// 	pgErr.Message,
			// 	pgErr.Detail,
			// 	pgErr.Where,
			// 	pgErr.Code,
			// 	pgErr.SQLState(),
			// )
			return nil, err
		}
		return nil, err
	}

	return canvases, nil
}

func (r *repository) GetByID(ctx context.Context, id int, userID int) (*canvas.CanvasBase, error) {
	q := `
		SELECT id, name, url, updated_at, text
		FROM canvas
		WHERE id = $1 AND user_id = $2
	`
	canvas := &canvas.CanvasBase{}

	row := r.client.QueryRow(ctx, q, id, userID)

	var nullUrl sql.NullString
	var nullText sql.NullString

	err := row.Scan(&canvas.CanvasID, &canvas.Name, &nullUrl, &canvas.UpdatedAt, &nullText)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
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

	if nullUrl.Valid {
		canvas.Url = nullUrl.String
	} else {
		canvas.Url = ""
	}
	if nullText.Valid {
		canvas.Text = nullText.String
	} else {
		canvas.Text = ""
	}

	return canvas, nil
}

func (r *repository) Delete(ctx context.Context, id int, userID int) error {
	query := `
		DELETE 
		FROM canvas
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.client.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *repository) ChangeParent(ctx context.Context, canvas_id int, new_parent_id int) error {
	q := `
		UPDATE canvas
		SET folder_id = $1
		WHERE id = $2
	`

	_, err := r.client.Exec(ctx, q, new_parent_id, canvas_id)
	if err != nil {
		return fmt.Errorf("error while updating canvas url")
	}
	return nil
}
