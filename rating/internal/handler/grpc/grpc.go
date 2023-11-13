package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/rating/internal/controller/rating"
	"movieexample.com/rating/pkg/model"
)

// Handler defined a movie metada gRPC handler.
type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

// New creates a new rating gRPC handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMetadata retuns a aggregated rating.
func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	v, err := h.ctrl.GetAggreatedRating(ctx, model.RecordType(req.RecordType), model.RecordID(req.RecordId))

	if err != nil && errors.Is(err, rating.ErrNotFound) {
		// return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())

	}
	return &gen.GetAggregatedRatingResponse{Rating: v}, nil
}

// PutRating writes a rating for a given record.
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	if err := h.ctrl.Put(ctx, model.RecordType(req.RecordType), model.RecordID(req.RecordId), &model.Rating{UserID: model.UserID(req.UserId), RatingValue: model.RatingValue(req.RatingValue)}); err != nil {
		return nil, err
	}
	return &gen.PutRatingResponse{}, nil
}
