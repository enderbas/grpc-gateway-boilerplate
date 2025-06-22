package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// Import your generated Go package
	gw "example.com/gateway/proto"
)

var (
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:50051", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// The gRPC-Gateway Mux
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := gw.RegisterGreeterHandlerFromEndpoint(ctx, gwmux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}

	// Create a new ServeMux for the main HTTP server
	mux := http.NewServeMux()

	// Mount the gRPC-Gateway proxy
	mux.Handle("/", gwmux)

	// Mount the handler for the swagger.json file
	mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "proto/greeter.swagger.json")
	})

	// Mount the Swagger UI file server
	fs := http.FileServer(http.Dir("swagger-ui"))
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", fs))

	log.Println("Starting gRPC-Gateway on http://0.0.0.0:8081")
	log.Println("Swagger UI available at http://localhost:8081/swagger-ui/")
	// Start HTTP server
	return http.ListenAndServe(":8081", allowCORS(mux))
}

// allowCORS allows Cross Origin Resourcing Sharing
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}


func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("Failed to start gRPC-Gateway: %v", err)
	}
}