package metadata

import (
	"context"
	"errors"

	repository "movieexample.com/metadata/internal/repository"
	models "movieexample.com/metadata/pkg/models"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*models.Metadata, error)
	Put(ctx context.Context, id string, metadata *models.Metadata) error
}

type Controller struct {
	repository metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repository: repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*models.Metadata, error) {
	res, err := c.repository.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	return res, nil
}
func (c *Controller) Put(ctx context.Context, metadata *models.Metadata) error {
	err := c.repository.Put(ctx, metadata.ID, metadata)

	if err != nil {
		return err
	}
	return nil
}
