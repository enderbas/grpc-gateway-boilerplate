# C++ gRPC Service with Go gRPC-Gateway

This project is a complete, working example of a high-performance C++ gRPC service that is exposed as a user-friendly RESTful JSON API using the Go-based gRPC-Gateway.

It is designed to be a boilerplate for building modern microservices where you want the internal performance of gRPC but need to provide a standard REST interface for web clients, legacy systems, or public consumption.

This repository was built and debugged to work specifically on a Debian-based Linux system (like Ubuntu) using the `apt` package manager for C++ dependencies.

## Architecture

The system runs as two separate processes that communicate with each other:

```
+--------------+      HTTP/1.1      +-------------------+      gRPC      +-----------------+
|              | <----------------> |                   | <------------> |                 |
| REST Client  |      (JSON)        |  Go Gateway Proxy |   (Protobuf)   | C++ gRPC Server |
| (e.g. curl)  |                    | (localhost:8081)  |                | (localhost:50051)|
|              | <----------------> |                   | <------------> |                 |
+--------------+      (JSON)        +-------------------+   (Protobuf)   +-----------------+
```

  - **C++ gRPC Server**: The core application logic written in C++. Handles requests with maximum performance.
  - **Go Gateway Proxy**: A lightweight reverse-proxy that translates incoming HTTP/JSON requests into gRPC requests, which are then forwarded to the C++ server.

## Prerequisites

Before you begin, you must have the following tools and libraries installed on your system. These instructions are tailored for **Ubuntu/Debian**.

### 1\. Core Build Tools

Install the essential C++ compiler toolchain and CMake. gRPC requires a modern version of CMake.

```bash
sudo apt update
sudo apt install build-essential cmake
```

### 2\. Go Language

Install a recent version of the Go programming language (1.20+ is recommended).
You can follow the official installation guide: [https://go.dev/doc/install](https://go.dev/doc/install)

After installing, ensure the Go binary path is added to your shell's `PATH`. Add this line to your `~/.bashrc` or `~/.zshrc`:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Then, reload your shell with `source ~/.bashrc` or open a new terminal.

### 3\. C++ gRPC & Protobuf Libraries

We will use the versions provided by `apt` for simplicity. These packages provide the compiled libraries and headers needed by the C++ server.

```bash
sudo apt install libgrpc-dev libprotobuf-dev protobuf-compiler grpc-proto
```

**Note**: The `grpc-proto` package is crucial as it provides the `grpc_cpp_plugin` executable needed by our build script.

### 4\. Go Code Generation Tools

Install the `protoc` plugins for Go. These are required to generate the Go client stubs and the gateway code from our `.proto` definition.

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

## Project Structure

```
.
├── CMakeLists.txt         # Build script for the C++ server
├── gateway/               # Contains all Go gateway code
│   ├── go.mod
│   └── main.go
├── proto/                 # The single source of truth for our API
│   ├── greeter.proto
│   └── google/api/...
└── server/                # Contains the C++ server implementation
    └── greeter_server.cc
```

## Step-by-Step Build and Run Instructions

### 1\. Clone the Repository

```bash
git clone <your-repository-url>
cd <your-repository-name>
```

### 2\. Generate Go Code for the Gateway

The C++ code generation is handled automatically by CMake, but we need to manually generate the Go code.

From the **root of the project directory**, run the following `protoc` command. This will create the necessary Go files inside the `gateway/proto/` directory.

```bash
# Create the target directory first
mkdir -p gateway/proto

# Run the generator
protoc -I ./proto \
    --go_out=./gateway/proto --go_opt=paths=source_relative \
    --go-grpc_out=./gateway/proto --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=./gateway/proto --grpc-gateway_opt=paths=source_relative \
    ./proto/greeter.proto
```

### 3\. Build the C++ gRPC Server

Use CMake to build the server executable.

```bash
# Create a build directory
mkdir build
cd build

# Configure the project and build
cmake ..
make
```

If the build is successful, you will find a `greeter_server` executable inside the `build` directory.

### 4\. Prepare the Go Gateway Module

Navigate to the `gateway` directory and run `go mod tidy`. This will scan the source code, find the required dependencies, and update the `go.mod` and `go.sum` files.

```bash
# Navigate from the build directory back to the root, then into gateway
cd ../gateway
go mod tidy
```

### 5\. Run the Application

The system requires two processes running in parallel. You will need **two separate terminals**.

**Terminal 1: Start the C++ gRPC Server**

```bash
# Navigate to your build directory
cd build

# Run the server
./greeter_server

# You should see:
# C++ gRPC Server listening on 0.0.0.0:50051
```

Leave this terminal running.

**Terminal 2: Start the Go Gateway Proxy**

```bash
# Navigate to the gateway directory
cd gateway

# Run the proxy
go run .

# You should see:
# Starting gRPC-Gateway on http://0.0.0.0:8081
```

Leave this terminal running as well.

### 6\. Test the System\!

With both servers running, open a **third terminal** and use `curl` to send an HTTP/JSON request to the gateway.

```bash
curl -X POST -k http://localhost:8081/v1/greeter/say_hello \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

**Response:**

You should receive a successful JSON response in your `curl` terminal:

```json
{
  "message": "Hello, World!"
}
```

Check the terminal running your C++ server. You will see the log output confirming it handled the request:

```
Received request from: World
```

Congratulations\! You have a fully functional gRPC + REST API service running.
