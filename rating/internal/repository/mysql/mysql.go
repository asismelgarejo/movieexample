package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
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

// Get retrieves all ratings for a give record
func (r *Repository) Get(ctx context.Context, recordType model.RecordType, recordID model.RecordID) ([]*model.Rating, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT user_id, value FROM ratings WHERE record_id = ? AND record_type = ?", recordID, recordType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*model.Rating
	for rows.Next() {
		var userID string
		var value int32
		if err := rows.Scan(&userID, &value); err != nil {
			return nil, err
		}
		res = append(res, &model.Rating{
			UserID:      model.UserID(userID),
			RatingValue: model.RatingValue(value),
		})
	}
	if len(res) == 0 {
		return nil, repository.ErrNotFound
	}
	return res, nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordType model.RecordType, recordID model.RecordID, rating *model.Rating) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO ratings (record_id, record_type, user_id, value) VALUES (?, ?, ?, ?)",
		recordID, recordType, rating.UserID, rating.RatingValue,
	)
	return err
}
