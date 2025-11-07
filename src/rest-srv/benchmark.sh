#!/bin/bash

# Benchmark script for REST API
# Usage: ./benchmark.sh [endpoint]
# Example: ./benchmark.sh "http://localhost:3000/person?id=1"

ENDPOINT="${1:-http://localhost:3000/person?id=1}"
THREADS="${2:-8}"
CONNECTIONS="${3:-400}"
DURATION="${4:-30s}"

echo "Running benchmark..."
echo "Endpoint: $ENDPOINT"
echo "Threads: $THREADS, Connections: $CONNECTIONS, Duration: $DURATION"
echo ""

wrk -t"$THREADS" -c"$CONNECTIONS" -d"$DURATION" "$ENDPOINT"

