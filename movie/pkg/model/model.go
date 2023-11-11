package model

import (
	metadata "movieexample.com/metadata/pkg/models"
)

type MovieDetails struct {
	Rating   float64            `json:"rating"`
	Metadata *metadata.Metadata `json:"metadata"`
}
