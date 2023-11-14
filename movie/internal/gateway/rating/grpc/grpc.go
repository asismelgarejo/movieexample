package grpc

import (
	"context"

	"movieexample.com/gen"
	"movieexample.com/internal/grpcutil"
	"movieexample.com/pkg/discovery"

	"movieexample.com/rating/pkg/model"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

func (g *Gateway) GetAggregatedRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID) (float64, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{RecordType: string(recordType), RecordId: string(recordId)})
	if err != nil {
		return 0, err
	}
	return resp.Rating, nil
}
func (g *Gateway) PutRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating model.Rating) error {
	conn, err := grpcutil.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := gen.NewRatingServiceClient(conn)

	_, err = client.PutRating(ctx, &gen.PutRatingRequest{UserId: string(rating.UserID), RecordId: string(rating.RecordID), RecordType: string(rating.RecordType), RatingValue: int32(rating.RatingValue)})
	if err != nil {
		return err
	}
	return nil
}
