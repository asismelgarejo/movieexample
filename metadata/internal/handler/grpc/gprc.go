package grpc

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/internal/controller/metadata"
	"movieexample.com/metadata/internal/repository"
	"movieexample.com/metadata/pkg/models"
)

// Handler defined a movie metada gRPC handler.
type Handler struct {
	gen.UnimplementedMetadataServiceServer
	ctrl *metadata.Controller
}

// New creates a new movie metadata gRPC handler.
func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMetadata retuns a movie metadata.
func (h *Handler) GetMetadata(ctx context.Context, req *gen.GetMetadataRequest) (*gen.GetMetadataResponse, error) {
	log.Printf("Starting GetMetadata")

	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	log.Printf("GetAggregatedRating:\treq.MovieId->%v\n", req.MovieId)

	m, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, status.Errorf(codes.InvalidArgument, "req nil or empty id")

	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetMetadataResponse{Metadata: models.MetadataToProto(m)}, nil
}

func (h *Handler) PutMetadata(ctx context.Context, req *gen.PutMetadataRequest) (*gen.PutMetadataResponse, error) {
	log.Printf("Starting PutMetadata")
	log.Printf("GetAggregatedRating:\treq.MovieId->%v\n", models.MetadataFromProto(req.Metadata))

	err := h.ctrl.Put(ctx, models.MetadataFromProto(req.Metadata))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.PutMetadataResponse{}, nil
}
