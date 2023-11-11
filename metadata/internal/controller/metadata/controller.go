package metadata

import (
	"context"
	"errors"

	repository "movieexample.com/metadata/internal/repository"
	models "movieexample.com/metadata/pkg/models"
)

type metadataRepository interface {
	Get(ctx context.Context, id string) (*models.Metadata, error)
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
	}
	return res, nil
}
