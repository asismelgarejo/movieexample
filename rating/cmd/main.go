package main

import (
	"fmt"
	"log"
	"net/http"

	controller "movieexample.com/rating/internal/controller/rating"
	handler "movieexample.com/rating/internal/handler/http"
	repository "movieexample.com/rating/internal/repository/memory"
)

func main() {
	fmt.Println("Starting the rating service")
	repo := repository.New()
	ctrl := controller.New(repo)
	h := handler.New(ctrl)

	http.Handle("/rating", http.HandlerFunc(h.Handle))
	log.Println("Server running on http://localhost:8082")
	log.Panic(http.ListenAndServe(":8082", nil))
}
