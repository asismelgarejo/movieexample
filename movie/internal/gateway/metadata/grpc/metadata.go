package http

import (
	"context"

	"movieexample.com/gen"
	"movieexample.com/internal/grpcutil"
	metadata "movieexample.com/metadata/pkg/models"
	"movieexample.com/pkg/discovery"
)

// Gateway defined a movie metadata HTTP gateway
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP movie gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	// func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Get gets a movie metadata using a movie id.

func (g *Gateway) Get(ctx context.Context, id string) (*metadata.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}
	return metadata.MetadataFromProto(resp.Metadata), nil
}
