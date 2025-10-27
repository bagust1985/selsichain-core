#!/bin/bash
echo "🚀 Starting SelsiChain Node on Railway..."
echo "🌐 PORT: $PORT"
echo "📁 Current directory: $(pwd)"
echo "📋 Files:"
ls -la

# Build binary
echo "🔨 Building binary..."
go build -o bin/selsichain ./cmd/selsichain/main.go

# Run the node
echo "🎯 Starting node..."
./bin/selsichain --p2p-port=$PORT --testnet