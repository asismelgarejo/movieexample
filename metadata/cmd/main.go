package main

import (
	"log"
	"net/http"

	controller "movieexample.com/metadata/internal/controller/metadata"
	handlers "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
)

func main() {
	memory := memory.New()
	ctrl := controller.New(memory)
	h := handlers.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	log.Println("Server running on http://localhost:8081")
	log.Panic(http.ListenAndServe(":8081", nil))
}
