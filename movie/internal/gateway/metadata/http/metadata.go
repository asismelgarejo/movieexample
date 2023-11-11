package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	metadata "movieexample.com/metadata/pkg/models"
	gateway "movieexample.com/movie/internal/gateway"
)

// Gateway defined a movie metadata HTTP gateway
type Gateway struct {
	addr string
}

// New creates a new HTTP movie gateway for a movie metadata service.
func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

// Get gets a movie metadata using a movie id.

func (g *Gateway) Get(ctx context.Context, id string) (*metadata.Metadata, error) {

	addr := "http://" + g.addr + "/metadata"
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx status code: %v", resp.StatusCode)
	}
	var metadata *metadata.Metadata

	if err = json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, err
	}
	log.Print("metadata", metadata)
	return metadata, nil
}
