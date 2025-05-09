package folder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/interfaces/folder"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	client *pgxpool.Pool
}

func NewRepository(client *pgxpool.Pool) folder.Repository {
	return &repository{client: client}
}

func (r *repository) GetFolder(ctx context.Context, folder_id int, user_id int) (*canvas.FolderBase, error) {
	folder := &canvas.FolderBase{}

	q := `
		SELECT id, name, parent_folder_id
		FROM folder
 		WHERE id = $1 AND user_id = $2
	`

	// print("\n", id, "\n")

	row := r.client.QueryRow(ctx, q, folder_id, user_id)
	err := row.Scan(&folder.ID, &folder.Name, &folder.Parent)
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

	return folder, nil

}

func (r *repository) Get(ctx context.Context, folder_id int, user_id int) (*canvas.FolderContent, error) {
	folder := &canvas.FolderBase{}

	q := `
		SELECT id, name, parent_folder_id, updated_at, created_at
		FROM folder
 		WHERE id = $1 AND user_id = $2
	`

	fmt.Println("Ваш user_id: ", user_id)

	if folder_id == 0 {
		id := 0
		folder.ID = &id
		folder.Name = "root"
	} else {
		row := r.client.QueryRow(ctx, q, folder_id, user_id)

		err := row.Scan(&folder.ID, &folder.Name, &folder.Parent, &folder.CreatedAt, &folder.UpdatedAt)
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
	}

	canvases := make([]canvas.CanvasBase, 0)

	q = `
		SELECT id, name, url, updated_at, text, folder_id, user_id, created_at
		FROM canvas
		WHERE folder_id = $1 AND user_id = $2
	`

	rows, err := r.client.Query(ctx, q, folder_id, user_id)
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

	print(canvases)

	folders := make([]canvas.FolderBase, 0)

	q = `
		SELECT id, name, parent_folder_id, updated_at, created_at
		FROM folder
		WHERE parent_folder_id = $1 AND user_id = $2
	`

	rows, err = r.client.Query(ctx, q, folder_id, user_id)
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

	for rows.Next() {
		var folder canvas.FolderBase

		err = rows.Scan(&folder.ID, &folder.Name, &folder.Parent, &folder.UpdatedAt, &folder.CreatedAt)
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
		folders = append(folders, folder)
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

	folderContent := canvas.FolderContent{
		Folder:   *folder,
		Folders:  folders,
		Canvases: canvases,
	}

	return &folderContent, nil

}

func (r *repository) Create(ctx context.Context, folder canvas.FolderBase, user_id int) (*int, error) {
	q := `
		INSERT INTO folder (name, parent_folder_id, user_id, updated_at, created_at) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id;
	`

	err := r.client.QueryRow(ctx, q, folder.Name, folder.Parent, user_id, folder.UpdatedAt, folder.CreatedAt).Scan(&folder.ID)
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

	return folder.ID, nil
}

func (r *repository) Delete(ctx context.Context, folder_id int) error {
	query := `
		DELETE 
		FROM folder
		WHERE id = $1
	`

	result, err := r.client.Exec(ctx, query, folder_id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
