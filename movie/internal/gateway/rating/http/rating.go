package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"movieexample.com/movie/internal/gateway"
	model "movieexample.com/rating/pkg/model"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

// GetAggregatedRating makes an http request to rating service for getting the aggregated rating for an specific record.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID) (float64, error) {

	addr := "http://" + g.addr + "/rating"
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return 0, err
	}
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("type", string(recordType))
	values.Add("id", string(recordId))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return 0, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return 0, fmt.Errorf("a non-2xx status code was returned: %d", resp.StatusCode)
	}
	var rating float64

	if err := json.NewDecoder(resp.Body).Decode(&rating); err != nil {
		return 0, err
	}
	return rating, nil
}

// PutRating makes an http request to rating service for adding a new rating
func (g *Gateway) PutRating(ctx context.Context, recordType model.RecordType, recordId model.RecordID, rating model.Rating) error {

	jsonData, _ := json.Marshal(rating)

	req, err := http.NewRequest(http.MethodPost, g.addr+"/rating", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	values := req.URL.Query()
	values.Add("type", string(recordType))
	values.Add("id", string(recordId))
	req.URL.RawQuery = values.Encode()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		return gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return fmt.Errorf("a non-2xx status code was returned: %d", resp.StatusCode)
	}

	return nil
}
