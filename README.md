# DSEQ - Distributed Sequence Generator

A distributed sequence generator built with Go and CometBFT, providing reliable and consistent sequence numbers across multiple nodes.

## Features

- Distributed sequence generation
- Multi-node support
- Load testing capabilities
- Health monitoring
- Metrics collection
- Tracing support
- Structured logging

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- CometBFT v0.38.0-rc3
- Make

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/dseq.git
cd dseq
```

2. Install dependencies:
```bash
go mod download
```

## Building

### Local Build
Build for your current platform:
```bash
make build
```

### Docker Build
Build multi-architecture Docker images:
```bash
make build-docker
```

Build local Docker image:
```bash
make build-docker-local
```

## Running

### Start Local Testnet
Start a 4-node testnet locally:
```bash
make start
```

Stop the testnet:
```bash
make stop
```

### Load Testing

Run load tests with different configurations:

```bash
# Quick test (10 requests, 3 concurrent)
make load-quick

# Medium test (100 requests, 9 concurrent)
make load-medium

# Heavy test (1000 requests, 30 concurrent)
make load-heavy

# Custom test
make load NODES=localhost:26657 REQUESTS=50 CONCURRENCY=5
```

### Monitoring

Monitor all nodes:
```bash
make read-all
```

Compare node sequence files:
```bash
make checksum
```

## Development

### Testing
Run tests with race detection and coverage:
```bash
make test
```

Generate test coverage report:
```bash
make test-coverage
```

### Code Quality
Format code:
```bash
make fmt
```

Run linters:
```bash
make lint
```

Run go vet:
```bash
make vet
```

### Cleanup
Clean build artifacts:
```bash
make clean
```

## Project Structure

```
.
├── app/                    # Application code
│   ├── config/            # Configuration
│   ├── errors/            # Error handling
│   ├── health/            # Health checks
│   ├── logging/           # Structured logging
│   ├── metrics/           # Metrics collection
│   ├── middleware/        # HTTP middleware
│   ├── testutil/          # Test utilities
│   └── tracing/           # Distributed tracing
├── build/                 # Build artifacts
├── cmd/                   # Command-line tools
├── networks/             # Network configurations
│   └── local/            # Local testnet setup
└── Makefile              # Build and development commands
```

## Available Make Targets

Run `make help` to see all available targets:

```bash
make help
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
