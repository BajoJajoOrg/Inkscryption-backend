package folder

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/BajoJajoOrg/Inkscryption-backend/canvas"
	"github.com/BajoJajoOrg/Inkscryption-backend/canvas/interfaces/folder"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/db/postgresql"
	"github.com/BajoJajoOrg/Inkscryption-backend/pkg/filter"
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

func (r *repository) Get(ctx context.Context, filterOptions filter.Options, folder_id int, user_id int) (*canvas.FolderContent, error) {
	folder := &canvas.FolderBase{}

	q := `
		SELECT id, name, parent_folder_id, updated_at, created_at
		FROM folder
 		WHERE id = $1 AND user_id = $2
	`

	// fmt.Println("Ваш user_id: ", user_id)

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

	// print(canvases)

	query, args, err = postgresql.BuildFolderQuery(filterOptions, folder_id, user_id)
	if err != nil {
		return nil, err
	}

	pgxArgs = make([]any, 0, len(args))
	for _, arg := range args {
		switch v := arg.(type) {
		// case time.Time:
		// 	pgxArgs = append(pgxArgs, v.Format(time.RFC3339))
		default:
			pgxArgs = append(pgxArgs, v)
		}
	}

	rows, err = r.client.Query(ctx, query, pgxArgs...)
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

	folders := make([]canvas.FolderBase, 0)

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

func (r *repository) ChangeFolderParent(ctx context.Context, folder_id int, new_parent_id int) error {
	q := `
		UPDATE folder
		SET parent_folder_id = $1
		WHERE id = $2
	`

	_, err := r.client.Exec(ctx, q, new_parent_id, folder_id)
	if err != nil {
		return fmt.Errorf("error while updating canvas url")
	}
	return nil
}

func (r *repository) ChangeCanvasParent(ctx context.Context, canvas_id int, new_parent_id int) error {
	q := `
		UPDATE canvas
		SET folder_id = $1
		WHERE id = $2
	`

	_, err := r.client.Exec(ctx, q, new_parent_id, canvas_id)
	if err != nil {
		return fmt.Errorf("error while updating canvas parent")
	}
	return nil
}

func (r *repository) Update(ctx context.Context, folder_id int, user_id int, name string) error {
	q := `
		UPDATE folder
		SET name = $1 
		WHERE id = $2 AND user_id = $3
	`

	_, err := r.client.Exec(ctx, q, name, folder_id, user_id)
	if err != nil {
		return fmt.Errorf("error while updating canvas name")
	}
	return nil
}
