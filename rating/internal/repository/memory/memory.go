package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	repository "movieexample.com/rating/internal/repository"
	"movieexample.com/rating/pkg/model"
)

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]*model.Rating
}

func GetCurrentDir() string {
	exePath, _ := os.Executable()
	return exePath
}

func New() *Repository {
	exePath := GetCurrentDir()
	filepath := filepath.Join(filepath.Dir(exePath), "/mocks/rating.json")
	recordsString, err := os.ReadFile(filepath)
	if err != nil {
		panic(fmt.Sprintf("Error when reading rating file: %v", err.Error()))
	}

	records := make(map[model.RecordType]map[model.RecordID][]*model.Rating)
	err = json.Unmarshal(recordsString, &records)
	if err != nil {
		panic(fmt.Sprintf("Error unmarshelling file content: %v", err.Error()))
	}
	fmt.Println(records)

	return &Repository{data: records}
}

// Get retrieves all ratings for a give record
func (r *Repository) Get(ctx context.Context, recordType model.RecordType, recordID model.RecordID) ([]*model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordID]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

// Put adds a new rating for a given record
func (r *Repository) Put(ctx context.Context, recordType model.RecordType, recordID model.RecordID, rating *model.Rating) error {
	if _, ok := r.data[recordType]; !ok {
		return repository.ErrNotFound
	}
	r.data[recordType][recordID] = append(r.data[recordType][recordID], rating)
	return nil
}
