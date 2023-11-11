package main

import (
	"log"
	"net/http"

	controller "movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/http"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/http"
	handler "movieexample.com/movie/internal/handler/http"
)

func main() {
	metadataGateway := metadatagateway.New("localhost:8081")
	ratingGateway := ratinggateway.New("localhost:8082")

	ctrl := controller.New(ratingGateway, metadataGateway)
	h := handler.New(ctrl)
	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))
	log.Println("Server running on http://localhost:8083")
	log.Panic(http.ListenAndServe(":8083", nil))
}
