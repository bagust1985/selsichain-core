#!/bin/bash
echo "🚀 Launching SelsiChain Testnet Cluster..."

# Node 1 (already running)
echo "📍 Node 1: port 7690 (already running)"

# Launch other nodes in background
./bin/selsichain --p2p-port=7691 --testnet &
./bin/selsichain --p2p-port=7692 --testnet &
./bin/selsichain --p2p-port=7693 --testnet &
./bin/selsichain --p2p-port=7694 --testnet &

echo "✅ Testnet cluster launched!"
echo "🔍 Monitor each node in separate terminals"
echo "🛑 Stop with: pkill -f selsichain"