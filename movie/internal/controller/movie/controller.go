package controller

import (
	"context"
	"errors"

	metadata "movieexample.com/metadata/pkg/models"
	"movieexample.com/movie/internal/gateway"
	model "movieexample.com/movie/pkg/model"
	ratingModel "movieexample.com/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateaway interface {
	GetAggregatedRating(ctx context.Context, recordType ratingModel.RecordType, recordId ratingModel.RecordID) (float64, error)
	PutRating(ctx context.Context, recordType ratingModel.RecordType, recordId ratingModel.RecordID, rating ratingModel.Rating) error
}
type metadataGateaway interface {
	Get(ctx context.Context, id string) (*metadata.Metadata, error)
}

type Controller struct {
	ratingGateaway   ratingGateaway
	metadataGateaway metadataGateaway
}

func New(ratingGateaway ratingGateaway, metadataGateaway metadataGateaway) *Controller {
	return &Controller{ratingGateaway, metadataGateaway}
}

// Get retrieves a movie details including movie metadata and its aggregated rating information.
func (c *Controller) Get(ctx context.Context, id string) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateaway.Get(ctx, id)
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		return nil, ErrNotFound
	}
	movieDetails := &model.MovieDetails{Metadata: metadata}
	rating, err := c.ratingGateaway.GetAggregatedRating(ctx, ratingModel.RecordTypeMovie, ratingModel.RecordID(id))
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		// return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	movieDetails.Rating = rating
	return movieDetails, nil
}
