package rating

import (
	"context"
	"errors"

	"movieexample.com/rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordType model.RecordType, recordId model.RecordID) ([]*model.Rating, error)
	Put(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating *model.Rating) error
}

type Controller struct {
	repo ratingRepository
}

func New(repo ratingRepository) *Controller {
	return &Controller{repo: repo}
}

// GetAggreatedRating returns the aggregated rating for a record or ErrNotFound if there is no ratings for it.
func (c *Controller) GetAggreatedRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordType, recordId)
	if err != nil {
		return 0, err
	}

	sum := float64(0)
	for _, rating := range ratings {
		sum += float64(rating.RatingValue)
	}
	return sum / float64(len(ratings)), nil
}

// Put writes a rating for a give record
func (c *Controller) Put(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating *model.Rating) error {
	return c.repo.Put(ctx, recordType, recordId, rating)
}
