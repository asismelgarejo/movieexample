package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gopkg.in/yaml.v2"

	"movieexample.com/gen"
	controller "movieexample.com/metadata/internal/controller/metadata"
	grpchandler "movieexample.com/metadata/internal/handler/grpc"
	repository "movieexample.com/metadata/internal/repository/mysql"
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
)

const serviceName = "metadata"

func main() {
	// var port int
	// flag.IntVar(&port, "port", 8081, "API handler port")
	// flag.Parse()

	f, err := os.Open("base.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	log.Printf("Starting the metadata service on port %v", cfg.APIConfig.Port)

	registry, err := consul.NewRegistry(fmt.Sprintf("go_consul:%v", cfg.APIConfig.PortConsul))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	// if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%v:%v", "localhost", cfg.APIConfig.Port)); err != nil {
	if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%v:%v", "container_metadata", cfg.APIConfig.Port)); err != nil {
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

	repo, err := repository.New()
	if err != nil {
		panic(err)
	}

	ctrl := controller.New(repo)

	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.APIConfig.Port)) // An empty string to listen on all available network interfaces.
	// lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", cfg.APIConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterMetadataServiceServer(srv, h)
	// http.Handle("/metadata", http.HandlerFunc(h.GetMetadata))

	// log.Panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
	log.Panic(srv.Serve(lis), nil)
}
