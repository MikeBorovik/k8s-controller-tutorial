# K8s Controller Tutorial

[![CI](https://github.com/MikeBorovik/k8s-controller-tutorial/actions/workflows/ci.yml/badge.svg)](https://github.com/MikeBorovik/k8s-controller-tutorial/actions/workflows/ci.yml)

A tutorial project demonstrating how to build a Kubernetes controller in Go. This project showcases the fundamental principles of working with the Kubernetes API and creating custom controllers.

## Features

- CLI interface built with Cobra
- Built-in HTTP server using FastHTTP
- Kubernetes API integration
- Flexible logging system with zerolog
- Ready-to-use Helm charts for Kubernetes deployment
- CI/CD pipeline with GitHub Actions

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

# Start the HTTP server
./k8s-controller-tutorial server --port 8080

# Configure logging level
./k8s-controller-tutorial --log-level debug server
```

### Available Commands

- `server` - start the HTTP server
- `list` - example command for working with lists
- `go_basic` - basic Go examples

## Project Structure

```
.
├── chart/                  # Helm charts for Kubernetes deployment
├── cmd/                    # CLI commands (using Cobra)
│   ├── root.go             # Root command
│   ├── server.go           # Server command
│   └── ...
├── Dockerfile              # Docker image build
├── go.mod                  # Go modules
├── go.sum                  # Go dependencies
├── LICENSE                 # License
├── main.go                 # Entry point
├── Makefile                # Build and test commands
└── README.md               # Documentation
```

## Development

### Testing

```bash
go test ./...
```

### Building Docker Image

```bash
docker build -t k8s-controller-tutorial:latest .
```

### CI/CD Pipeline

This project uses GitHub Actions for continuous integration and delivery:

- Automatically builds and tests the code on push and pull requests
- Builds and publishes Docker images to GitHub Container Registry
- Packages Helm charts
- Runs security scanning with Trivy

The workflow is defined in `.github/workflows/ci.yml`.

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Author

MikeBorovik