# SpiceDB Go Library Examples

This directory contains example applications demonstrating how to use the SpiceDB Go library.

## HTTP Service Example

A complete HTTP service that wraps the SpiceDB library:

```bash
cd http-service
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
export SPICEDB_INSECURE=true
go run main.go
```

Then test it:

```bash
curl -X POST http://localhost:8080/check \
  -H "Content-Type: application/json" \
  -d '{
    "subject": "user:alice",
    "resource": "document:budget-2026",
    "permission": "view"
  }'
```

## Basic Examples

Comprehensive examples demonstrating:
- Basic permission checks
- Builder pattern
- Context/caveat usage
- Different consistency levels
- Batch checks
- Error handling
- Caching configuration
- Timeout management

```bash
cd basic
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
go run main.go
```

## Getting Started

1. Set up a local SpiceDB instance:

```bash
docker run --rm -p 50051:50051 authzed/spicedb serve \
  --grpc-preshared-key "devkey" \
  --datastore-engine memory
```

2. Run any of the examples with environment variables:

```bash
export SPICEDB_ADDR=localhost:50051
export SPICEDB_TOKEN=devkey
export SPICEDB_INSECURE=true
go run ./examples/basic/main.go
```

## Environment Variables

- `SPICEDB_ADDR`: SpiceDB server address (default: `localhost:50051`)
- `SPICEDB_TOKEN`: Pre-shared key for authentication
- `SPICEDB_INSECURE`: Use insecure connection (default: `true` for development)
- `SPICEDB_CONSISTENCY`: Default consistency level (default: `minimize_latency`)
- `SPICEDB_TIMEOUT`: Request timeout in seconds (default: `3s`)

## Library Documentation

For detailed library documentation, see [../spicedb/README.md](../spicedb/README.md)

