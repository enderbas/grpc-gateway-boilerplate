package main

import (
    "context"
    "flag"
    "log"
    "net/http"

    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"

    // This is the generated Go package from our .proto file
    gw "example.com/gateway/proto"
)

var (
    // command-line options:
    // gRPC server endpoint
    grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func run() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Register gRPC server endpoint
    // Note: Make sure the gRPC server is running properly and accessible
    mux := runtime.NewServeMux()
    opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
    
    // This is the key part: it registers the handler for our Greeter service
    err := gw.RegisterGreeterHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
    if err != nil {
        return err
    }

    log.Println("Starting gRPC-Gateway on http://0.0.0.0:8081")
    // Start HTTP server (and proxy calls to gRPC server)
    return http.ListenAndServe(":8081", mux)
}

func main() {
    flag.Parse()

    if err := run(); err != nil {
        log.Fatalf("Failed to start gRPC-Gateway: %v", err)
    }
}