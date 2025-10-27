package main

import (
    "os"
)

func main() {
    // Delegate to the actual main in cmd/selsichain
    os.Args = append([]string{os.Args[0], "--p2p-port=7690", "--testnet"}, os.Args[1:]...)
    
    // Import and run the actual main
    // This will be replaced by the actual implementation
    println("🚀 SelsiChain Node Starting...")
    println("📁 Redirecting to cmd/selsichain/main.go...")
    
    // For now, we'll copy the main logic here temporarily
    startFullNode("7690", true)
}
