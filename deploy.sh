#!/bin/bash
# deploy.sh

echo "🚀 Deploying SelsiChain Testnet..."

# Build binary
go build -o selsichain ./cmd/selsichain/main.go

# Create systemd service
sudo tee /etc/systemd/system/selsichain.service > /dev/null <<EOF
[Unit]
Description=SelsiChain Node
After=network.target

[Service]
Type=simple
User=ubuntu
WorkingDirectory=/home/ubuntu/selsichain-core
ExecStart=/home/ubuntu/selsichain-core/selsichain --p2p-port=7690 --testnet
Restart=always

[Install]
WantedBy=multi-user.target
EOF

# Start service
sudo systemctl daemon-reload
sudo systemctl enable selsichain
sudo systemctl start selsichain

echo "✅ SelsiChain deployed and running!"