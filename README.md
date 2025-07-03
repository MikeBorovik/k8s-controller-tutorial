# K8s Controller Tutorial

[![CI](https://github.com/MikeBorovik/k8s-controller-tutorial/actions/workflows/ci.yml/badge.svg)](https://github.com/MikeBorovik/k8s-controller-tutorial/actions/workflows/ci.yml)

A tutorial project demonstrating how to build a Kubernetes controller in Go. This project showcases the fundamental principles of working with the Kubernetes API and creating custom controllers.

## Features

- CLI interface built with Cobra
- Built-in HTTP server using FastHTTP
- Kubernetes API integration with client-go
- Deployment controller for reconciling Deployment resources
- Deployment informers for real-time monitoring
- Prometheus metrics for monitoring controller performance
- Leader election for high availability in multi-replica deployments
- Flexible logging system with zerolog
- Ready-to-use Helm charts for Kubernetes deployment
- CI/CD pipeline with GitHub Actions
- Comprehensive test suite

## Requirements

- Go 1.24+
- Kubernetes cluster (for testing the controller)
- Docker (for building images)
- Helm (for Kubernetes deployment)

## Installation

### From Source

```bash
git clone https://github.com/MikeBorovik/k8s-controller-tutorial.git
cd k8s-controller-tutorial
go build -o k8s-controller-tutorial
```

### Using Docker

```bash
docker build -t k8s-controller-tutorial .
docker run -p 8080:8080 k8s-controller-tutorial
```

### In Kubernetes with Helm

```bash
helm install k8s-controller ./chart/app
```

## Usage

### Running the CLI

```bash
# Show help
./k8s-controller-tutorial --help

# Start the HTTP server and deployment controller
./k8s-controller-tutorial server --port 8080 --kubeconfig ~/.kube/config

# Use in-cluster configuration when running in Kubernetes
./k8s-controller-tutorial server --in-cluster

# Configure metrics and leader election
./k8s-controller-tutorial server --metrics-port 8081 --enable-leader-election=true

# List deployments in the default namespace
./k8s-controller-tutorial list --kubeconfig ~/.kube/config

# Configure logging level
./k8s-controller-tutorial --log-level debug server
./k8s-controller-tutorial --log-level trace --log-format console server
```

### Available Commands

- `server` - start the HTTP server, deployment informer, and deployment controller
- `list` - list deployments in the default namespace

## Deployment Controller

The deployment controller is implemented using the `controller-runtime` library. It watches for changes to `Deployment` resources in the Kubernetes cluster and reconciles them. This is useful for implementing custom logic for managing deployments.

### Key Features
- Watches for `Deployment` resource changes
- Logs reconciliation events
- Uses `controller-runtime` for efficient resource management
- Exposes Prometheus metrics for monitoring
- Supports leader election for high availability

The controller is started automatically when you run the `server` command.

## Metrics

The controller exposes Prometheus metrics on a dedicated port (default: 8081). These metrics include:

- Standard controller-runtime metrics (reconciliation counts, durations, etc.)
- Custom metrics specific to deployment processing
- Go runtime metrics (memory usage, goroutines, etc.)

You can access metrics by navigating to `http://localhost:8081/metrics` when the server is running.

## Leader Election

When running multiple replicas of the controller in a Kubernetes cluster, leader election ensures that only one instance is actively reconciling resources at a time. This prevents conflicts and duplicate processing.

Leader election is enabled by default and can be configured with the `--enable-leader-election` flag. The leader election mechanism uses Kubernetes ConfigMaps to coordinate between replicas.

### Configuration

```bash
# Enable leader election (default)
./k8s-controller-tutorial server --enable-leader-election=true

# Disable leader election
./k8s-controller-tutorial server --enable-leader-election=false
```

## Project Structure

```
.
├── chart/                           # Helm charts for Kubernetes deployment
├── cmd/                             # CLI commands (using Cobra)
│   ├── root.go                      # Root command
│   ├── server.go                    # Server command with FastHTTP
│   ├── list.go                      # List command for K8s resources
│   └── ...
├── pkg/                             # Package code
│   ├── informer/                    # Kubernetes informers
│   │   └── informer.go              # Deployment informer implementation
│   └── ctrl/                        # Deployment controller
│       └── deployment_controller.go # Deployment controller implementation
├── Dockerfile                       # Docker image build
├── go.mod                           # Go modules
├── go.sum                           # Go dependencies
├── LICENSE                          # License
├── main.go                          # Entry point
├── Makefile                         # Build and test commands
└── README.md                        # Documentation
```

## Development

### Testing

```bash
# Run all tests using the test environment (recommended)
make test

# Run tests directly (without test environment)
go test ./...

# Skip informer tests in short mode
go test -short ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./cmd
go test ./pkg/informer
```

The `make test` command sets up the proper test environment using `setup-envtest` and handles dependencies automatically, making it the preferred way to run tests. It also generates JUnit XML reports for CI integration.

For code coverage reports:

```bash
# Generate coverage reports
make test-coverage
```

This will create both a standard Go coverage report and a Cobertura XML report for CI tools.

### Development Commands

```bash
# Format code
make format

# Run linter
make lint

# Build the binary
make build

# Run the application
make run

# Clean up build artifacts
make clean
```

### Building Docker Image

```bash
# Using docker directly
docker build -t k8s-controller-tutorial:latest .

# Using Makefile
make docker-build
```

### CI/CD Pipeline

This project uses GitHub Actions for continuous integration and delivery:

- Automatically builds and tests the code on push and pull requests
- Builds and publishes Docker images to GitHub Container Registry
- Packages Helm charts
- Runs security scanning with Trivy

The workflow is defined in `.github/workflows/ci.yml`.

## Key Components

### FastHTTP Server

The project includes a high-performance HTTP server using the FastHTTP library, which is significantly faster than the standard Go HTTP server.

### Kubernetes Informers

The controller uses Kubernetes informers to efficiently watch for changes to Deployment resources. This allows real-time monitoring without constant polling of the API server.

### Metrics Server

The controller includes a Prometheus metrics server that exposes metrics about controller performance and resource usage. These metrics can be scraped by Prometheus and visualized in dashboards.

### Leader Election

For high availability deployments, the controller supports leader election to ensure that only one instance is actively reconciling resources at a time. This prevents conflicts and duplicate processing when running multiple replicas.

### Structured Logging

Zerolog provides structured, JSON-formatted logs with configurable log levels (trace, debug, info, warn, error) and output formats.

### Test Environment

The project uses `setup-envtest` from the Kubernetes controller-runtime project to create a proper testing environment for Kubernetes API interactions. This allows tests to run against a real API server without needing an actual Kubernetes cluster.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Author

MikeBorovik