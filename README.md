# Project Setup Guide

## Prerequisites

Ensure that the MySQL container (`movieexample_db`) and Consul container (`go_consul`) are running before proceeding with the following steps.

## Step 1: Compile the Application

Compile the `cmd/main` package for each service using the following command:

### For Windows

```bash
go run -o app.exe ./cmd
```

### For PRODUCTION (Linux) on Windows

```bash
$Env:GOOS = 'linux'
go build -o exe\main .\cmd\
```

## Step 2: Run the Application

Execute the compiled application with the appropriate environment settings based on the deployment mode:

### Development (DEV)

```bash
./app.exe
```

## Production (PROD)

### Create the images

```docker build -t metadata .
docker build -t rating .
docker build -t movie .
```

### Create the containers

```# metadata
docker run --network go_micro_net --name container_metadata -p 8086:8086  metadata

# rating
docker run --network go_micro_net --name container_rating -p 8082:8082 rating

# movie
docker run --network go_micro_net  --name container_movie -p 8083:8083 -p 8084:8084 movie
```

Ensure you follow these steps in sequence to successfully set up and run the application. If you encounter any issues, please check the container status, compilation output, and ensure that the required dependencies are installed.
