package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"
	"movieexample.com/gen"
	controller "movieexample.com/movie/internal/controller/movie"
	metadatagateway "movieexample.com/movie/internal/gateway/metadata/grpc"
	ratinggateway "movieexample.com/movie/internal/gateway/rating/grpc"
	grpchandler "movieexample.com/movie/internal/handler/grpc"
	httphandler "movieexample.com/movie/internal/handler/http"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "movie"

func main() {
	// var port int
	// flag.IntVar(&port, "port", 8083, "API handler port")
	// flag.Parse()
	// log.Printf("Starting the movie service on port %v", port)

	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	log.Printf("Starting the movie service on port %v", cfg.APIConfig.Port)

	registry, err := consul.NewRegistry(fmt.Sprintf("go_consul:%v", cfg.APIConfig.PortConsul))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%v:%v", "localhost", cfg.APIConfig.Port)); err != nil {
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

	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)
	ctrl := controller.New(ratingGateway, metadataGateway)

	go func() {
		grpcHandler := grpchandler.New(ctrl)
		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.APIConfig.Port))
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		srv := grpc.NewServer()
		reflection.Register(srv)
		gen.RegisterMovieServiceServer(srv, grpcHandler)
		log.Printf("tcp server started %v", cfg.APIConfig.Port)
		log.Panic(srv.Serve(lis))
	}()

	httpHandler := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(httpHandler.GetMovieDetails))
	log.Panic(http.ListenAndServe(fmt.Sprintf(":%v", cfg.APIConfig.Port), nil))
}
