package mysql

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func New() (*Repository, error) {
	db, err := sql.Open("mysql", "root:password@/movieexample")
	if err != nil {
		return nil, err
	}
	return &Repository{db: db}, nil
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(ctx context.Context, id string) (*models.Metadata, error) {
	var title, description, director string
	row := r.db.QueryRowContext(ctx, "SELECT title,description, director FROM movies WHERE id = ?", id)
	if err := row.Scan(&title, &description, &director); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &models.Metadata{
		ID:          id,
		Title:       title,
		Description: description,
		Director:    director,
	}, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(ctx context.Context, id string, metadata *models.Metadata) error {
	log.Print("asdasdasd", id, metadata)
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO movies (id,title, description, director) VALUES (?, ?, ?, ?)", id, metadata.Title, metadata.Description, metadata.Director)
	return err
}
