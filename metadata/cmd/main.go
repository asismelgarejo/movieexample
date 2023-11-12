package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	controller "movieexample.com/metadata/internal/controller/metadata"
	handlers "movieexample.com/metadata/internal/handler/http"
	"movieexample.com/metadata/internal/repository/memory"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "API handler port")
	flag.Parse()
	log.Printf("Starting the metadata service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%v:%d", "localhost", port)); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	memory := memory.New()
	ctrl := controller.New(memory)
	h := handlers.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
