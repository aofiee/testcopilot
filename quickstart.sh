#!/bin/bash

# Quick Start Script for HMAC Service
# This script builds and runs the service, then shows some examples

set -e

echo "🚀 HMAC Service Quick Start"
echo "=========================="
echo

# Build the application
echo "📦 Building the application..."
go build -o bin/hmac-service ./cmd/server
echo "✅ Build complete!"
echo

# Start the service in the background
echo "🔧 Starting the service..."
PORT=8080 ./bin/hmac-service &
SERVER_PID=$!

# Wait for the service to start
echo "⏳ Waiting for service to start..."
for i in {1..10}; do
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        break
    fi
    sleep 1
done

# Check if service is running
if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "❌ Service failed to start"
    kill $SERVER_PID 2>/dev/null || true
    exit 1
fi

echo "✅ Service is running on http://localhost:8080"
echo

# Run examples
echo "📋 Running API examples..."
echo
./examples/curl_examples.sh

# Cleanup
echo
echo "🧹 Stopping the service..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo "✅ Done! You can now:"
echo "  - Run 'make run' to start the service"
echo "  - Run 'make test' to run tests"
echo "  - Run 'make docker-build' to build Docker image"
echo "  - See Makefile for more options"