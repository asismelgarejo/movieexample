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
	"movieexample.com/pkg/discovery"
	"movieexample.com/pkg/discovery/consul"
	controller "movieexample.com/rating/internal/controller/rating"
	grpchandler "movieexample.com/rating/internal/handler/grpc"
	repository "movieexample.com/rating/internal/repository/mysql"
)

const serviceName = "rating"

func main() {
	baseFile := "base_dev.yaml"
	if os.Getenv("mode") == "production" {
		baseFile = "base_prod.yaml"
	}

	f, err := os.Open(fmt.Sprintf("./configs/%v", baseFile))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var cfg serviceConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		panic(err)
	}
	log.Printf("Starting the rating service on port %v", cfg.GRPCConfig.Port)

	registry, err := consul.NewRegistry(fmt.Sprintf("%v:%v", cfg.ConsulConfig.Addr, cfg.ConsulConfig.Port))
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err = registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%v:%v", cfg.GRPCConfig.Addr, cfg.GRPCConfig.Port)); err != nil {
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

	repo, err := repository.New(cfg.DBConfig.StrConn)
	if err != nil {
		panic(err)
	}
	ctrl := controller.New(repo, nil)
	// -----------------------
	h := grpchandler.New(ctrl)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCConfig.Port)) // An empty string to listen on all available network interfaces.
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	reflection.Register(srv)
	gen.RegisterRatingServiceServer(srv, h)
	// -----------------------

	log.Panic(srv.Serve(lis), nil)
}
