# C++ gRPC Service with Go gRPC-Gateway and Swagger UI

This project is a complete, working example of a high-performance C++ gRPC service that is exposed as a user-friendly RESTful JSON API using the Go-based gRPC-Gateway.

It is designed to be a boilerplate for building modern microservices where you want the internal performance of gRPC but need to provide a standard REST interface for web clients, legacy systems, or public consumption.
It now includes an auto-generated, interactive **Swagger UI** for easy API exploration and testing, built directly from the `.proto` service definition.

This repository was built and debugged to work specifically on a Debian-based Linux system (like Ubuntu) using the `apt` package manager for C++ dependencies.

## Architecture

The system runs as two separate processes that communicate with each other, with the gateway now also serving the API documentation.

```
+--------------+      HTTP/1.1      +-------------------+      gRPC      +-----------------+
|              | <----------------> |                   | <------------> |                 |
| REST Client  |      (JSON)        |  Go Gateway Proxy |   (Protobuf)   | C++ gRPC Server |
|(curl/Browser)|                    | (localhost:8081)  |                |(localhost:50051)|
|              | <----------------> |                   | <------------> |                 |
+--------------+     (JSON/HTML)    +-------------------+   (Protobuf)   +-----------------+
                                    |
                                    +--> Serves /swagger-ui/
                                    +--> Serves /swagger.json
```

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

Install the `protoc` plugins for Go. These are required to generate the Go client stubs, the gateway code, and the OpenAPI v2 specification.

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

## Step-by-Step Setup and Build

### 1\. Clone the Repository

### 2\. Download Swagger UI (One-Time Setup)

The Swagger UI is a set of static HTML, CSS, and JavaScript files. We need to download them once and place them in the project.

1.  Go to the [Swagger UI releases page](https://github.com/swagger-api/swagger-ui/releases).
2.  Download the latest release (e.g., `swagger-ui-vX.Y.Z.zip`).
3.  Unzip the package. Inside, you will find a `dist` directory.
4.  Create a new directory in our project: `gateway/swagger-ui`.
5.  Copy the **entire contents** of the `dist` directory into your new `gateway/swagger-ui` directory.
6.  **Important**: Edit `gateway/swagger-ui/index.html` and change the `url` field to point to your local definition: `url: "/swagger.json"`.

### 3\. Generate Go Code & OpenAPI Spec

From the **root of the project directory**, run the following `protoc` command. This will generate all the necessary Go files **and** the `greeter.swagger.json` API definition inside the `gateway/proto/` directory.

```bash
# Create the target directory first
mkdir -p gateway/proto

# Run the full generator command
protoc -I ./proto \
    --go_out=./gateway/proto --go_opt=paths=source_relative \
    --go-grpc_out=./gateway/proto --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=./gateway/proto --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=./gateway/proto \
    ./proto/greeter.proto
```

### 4\. Build the C++ gRPC Server

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

### 5\. Prepare the Go Gateway Module

Navigate to the `gateway` directory and run `go mod tidy`. This will scan the source code, find the required dependencies, and update the `go.mod` and `go.sum` files.

```bash
# Navigate from the build directory back to the root, then into gateway
cd ../gateway
go mod tidy
```

## Running the Application

The system requires two processes running in parallel. You will need **two separate terminals**.

#### Terminal 1: Start the C++ gRPC Server

```bash
# Navigate to your build directory
cd build

# Run the server
./greeter_server

# You should see:
# C++ gRPC Server listening on 0.0.0.0:50051
```

Leave this terminal running.

#### Terminal 2: Start the Go Gateway Proxy

```bash
# Navigate to the gateway directory
cd gateway

# Run the proxy
go run .

# You should see the new output:
# Starting gRPC-Gateway on http://0.0.0.0:8081
# Swagger UI available at http://localhost:8081/swagger-ui/
```

Leave this terminal running as well.

## Testing the System

You can now test the API in two ways.

### 1\. Testing with Swagger UI (Recommended)

1.  Open your web browser and navigate to: **[http://localhost:8081/swagger-ui/](http://localhost:8081/swagger-ui/)**
2.  You should see the "Greeter Service API" documentation page.
3.  Expand the `/v1/greeter/say_hello` endpoint.
4.  Click the **"Try it out"** button.
5.  Edit the example JSON request body.
6.  Click **"Execute"**.

You will see the live response from your C++ server directly in the browser\!

### 2\. Testing with curl

You can still use `curl` from a third terminal to test the API endpoint directly.

```bash
curl -X POST -k http://localhost:8081/v1/greeter/say_hello \
  -H "Content-Type: application/json" \
  -d '{"name": "World"}'
```

Check the terminal running your C++ server. You will see the log output (`Received request from: World`) confirming it handled the request.