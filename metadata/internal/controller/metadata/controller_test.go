package metadata

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	gen "movieexample.com/gen/mock/metadata/repository"
	"movieexample.com/metadata/internal/repository"
	models "movieexample.com/metadata/pkg/models"
)

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"movieexample.com/metadata/internal/repository"
// 	models "movieexample.com/metadata/pkg/models"
// )

// type fakeMetadataRepository struct {
// 	returnRes *models.Metadata
// 	returnErr error
// }

// func (r *fakeMetadataRepository) SetReturnValues(res *models.Metadata, err error) {
// 	r.returnErr = err
// 	r.returnRes = res
// }

// func (r *fakeMetadataRepository) Get(ctx context.Context, id string) (*models.Metadata, error) {
// 	return r.returnRes, r.returnErr
// }
// func (r *fakeMetadataRepository) Put(ctx context.Context, id string, metadata *models.Metadata) error {
// 	return r.returnErr
// }

// func TestController(t *testing.T) {
// 	m := &fakeMetadataRepository{}
// 	m.SetReturnValues(nil, errors.New("should"))
// 	// m.SetReturnValues(nil, repository.ErrNotFound)
// 	ctrl := New(m)
// 	_, err := ctrl.Get(context.Background(), "1")
// 	if err != nil && !errors.Is(err, repository.ErrNotFound) {
// 		t.Errorf("ErrorNotFound expected but got %v", err.Error())
// 	} else if err == nil {
// 		t.Errorf("Error expected but got nil %v", err)
// 	}
// }

func TestController(t *testing.T) {
	tests := []struct {
		name       string
		expRepoRes *models.Metadata
		expRepoErr error
		wantRes    *models.Metadata
		wantErr    error
	}{
		{
			name:       "not found",
			expRepoErr: repository.ErrNotFound,
			wantErr:    ErrNotFound,
		},
		{
			name:       "unexpected error",
			expRepoErr: errors.New("unexpected error"),
			wantErr:    errors.New("unexpected error"),
		},
		{
			name:       "sucess",
			expRepoRes: &models.Metadata{},
			wantRes:    &models.Metadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repoMock := gen.NewMockmetadataRepository(ctrl)
			c := New(repoMock)

			ctx := context.Background()
			id := "id"

			repoMock.EXPECT().Get(ctx, id).Return(tt.expRepoRes, tt.expRepoErr)

			res, err := c.Get(ctx, id)
			assert.Equal(t, tt.expRepoErr, err, tt.name)
			assert.Equal(t, tt.expRepoRes, res, tt.name)
		})
	}
}
