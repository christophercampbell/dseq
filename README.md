# DSEQ - Distributed Byte Sequencer

A distributed byte sequencer built with Go and CometBFT that produces a uniform sequence of bytes from multiple uncoordinated producers. The system ensures that all nodes maintain the same sequence of bytes, even when they receive different input transactions.

## Features

- Distributed byte sequence coordination
- Byzantine Fault Tolerance (BFT) consensus
- Multi-node support with independent inputs
- Consistent byte sequence across all nodes
- Load testing capabilities
- Health monitoring
- Metrics collection
- Tracing support
- Structured logging

## How It Works

DSEQ uses CometBFT's BFT consensus to ensure that all nodes maintain the same sequence of bytes, regardless of the order or source of incoming transactions. This is achieved through:

1. **Transaction Broadcasting**: Nodes broadcast their transactions to the network
2. **Consensus**: CometBFT's BFT consensus ensures all nodes agree on the transaction order
3. **Sequence Generation**: Each node processes the same sequence of transactions in the same order
4. **Consistency**: All nodes produce identical byte sequences despite receiving different inputs

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

### Testing the Sequencer

1. Send transactions to different nodes:
```bash
# Send to node 1
curl -s 'localhost:26657/broadcast_tx_commit?tx="0xDEADBEEF"'

# Send to node 2
curl -s 'localhost:26660/broadcast_tx_commit?tx="0xCAFEBABE"'
```

2. Verify sequence consistency:
```bash
# Compare sequence files across nodes
make checksum
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

Monitor all nodes to verify sequence consistency:
```bash
make read-all
```

Compare node sequence files to ensure they match:
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
