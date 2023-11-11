package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	repository "movieexample.com/metadata/internal/repository"
	models "movieexample.com/metadata/pkg/models"
)

type Repository struct {
	sync.RWMutex
	data map[string]*models.Metadata
}

func GetDirectory() (string, error) {
	exePath, _ := os.Executable()
	return exePath, nil
}

func New() *Repository {
	exePath, _ := GetDirectory()
	moviesPath := filepath.Join(filepath.Dir(exePath), "mocks/movies.json")
	moviesString, err := os.ReadFile(moviesPath)

	if err != nil {
		panic(fmt.Sprintf("Error reading movies.json: %v", err.Error()))
	}

	// Unmarshal JSON data into a slice of Metadata
	var movies []models.Metadata
	err = json.Unmarshal(moviesString, &movies)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshaling JSON: %v", err.Error()))
	}
	metadata := make(map[string]*models.Metadata)
	for _, movie := range movies {
		currentMovie := movie
		metadata[movie.ID] = &currentMovie
	}

	return &Repository{data: metadata}
}
func (r *Repository) Get(_ context.Context, id string) (*models.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	metadata, ok := r.data[id]

	if !ok {
		return nil, repository.ErrNotFound
	}
	return metadata, nil
}
