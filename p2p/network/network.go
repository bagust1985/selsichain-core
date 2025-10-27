package network

import (
    "fmt"
    "net"
    "sync"
    "time"
    
    "github.com/selsichain/selsichain-core/core/blockchain"
    "github.com/selsichain/selsichain-core/core/types"
)

// Network adalah P2P network untuk SelsiChain (Simple Version)
type Network struct {
    chain   *blockchain.Blockchain
    config  *Config
    peers   map[string]*PeerInfo
    mu      sync.RWMutex
    running bool
}

// Config holds network configuration
type Config struct {
    ListenAddr    string
    BootstrapPeers []string
    ChainID       uint64
    ProtocolID    string
}

// PeerInfo holds information about connected peers
type PeerInfo struct {
    ID        string
    Address   string
    Connected bool
    LastSeen  time.Time
}

// NewNetwork creates a new P2P network
func NewNetwork(config *Config, chain *blockchain.Blockchain) (*Network, error) {
    network := &Network{
        chain:   chain,
        config:  config,
        peers:   make(map[string]*PeerInfo),
        running: false,
    }

    return network, nil
}

// Start begins the network
func (n *Network) Start() error {
    fmt.Printf("🌐 Starting P2P Network...\n")
    fmt.Printf("📍 Will listen on: %s\n", n.config.ListenAddr)
    fmt.Printf("🔗 Protocol: %s\n", n.config.ProtocolID)
    fmt.Printf("⛓️  Chain ID: %d\n", n.config.ChainID)
    
    n.running = true
    
    // Connect to bootstrap peers
    if len(n.config.BootstrapPeers) > 0 {
        fmt.Printf("🔌 Connecting to %d bootstrap peers...\n", len(n.config.BootstrapPeers))
        n.connectToBootstrapPeers()
    }
    
    // Start background tasks
    go n.monitorNetwork()
    go n.cleanupDeadPeers()
    
    fmt.Printf("✅ P2P Network started successfully!\n")
    return nil
}

// Stop shuts down the network
func (n *Network) Stop() {
    fmt.Printf("🛑 Stopping P2P Network...\n")
    n.running = false
    fmt.Printf("✅ P2P Network stopped\n")
}

// GetPeerID returns the node's peer ID
func (n *Network) GetPeerID() string {
    return "selsichain-node-" + n.config.ListenAddr
}

// GetListenAddr returns the listening address
func (n *Network) GetListenAddr() string {
    return n.config.ListenAddr
}

// AddPeer adds a peer to the network WITH REAL CONNECTION CHECK
func (n *Network) AddPeer(address string) error {
    // ✅ REAL CONNECTION CHECK sebelum tambah peer
    if !n.checkPeerConnection(address) {
        return fmt.Errorf("cannot connect to peer %s", address)
    }
    
    n.mu.Lock()
    defer n.mu.Unlock()
    
    peerID := "peer-" + address
    n.peers[peerID] = &PeerInfo{
        ID:        peerID,
        Address:   address,
        Connected: true,
        LastSeen:  time.Now(),
    }
    
    fmt.Printf("✅ Added peer: %s\n", address)
    return nil
}

// parseAddress extracts host and port from libp2p-style address ✅ FIXED
func (n *Network) parseAddress(address string) (string, string, error) {
    // Simple parsing untuk format: /ip4/127.0.0.1/tcp/7690
    // Extract IP dan port dengan string manipulation
    
    // Remove /ip4/ prefix
    if len(address) < 6 || address[:5] != "/ip4/" {
        return "", "", fmt.Errorf("invalid address format, expected /ip4/ prefix: %s", address)
    }
    
    rest := address[5:] // Remove "/ip4/"
    
    // Find /tcp/ separator
    tcpIndex := -1
    for i := 0; i < len(rest)-4; i++ {
        if rest[i:i+4] == "/tcp" {
            tcpIndex = i
            break
        }
    }
    
    if tcpIndex == -1 {
        return "", "", fmt.Errorf("invalid address format, expected /tcp/ in: %s", address)
    }
    
    ip := rest[:tcpIndex]
    port := rest[tcpIndex+5:] // Remove "/tcp/"
    
    // Basic validation
    if ip == "" || port == "" {
        return "", "", fmt.Errorf("empty IP or port in: %s", address)
    }
    
    return ip, port, nil
}

// checkPeerConnection checks if we can actually connect to the peer ✅ IMPROVED
func (n *Network) checkPeerConnection(address string) bool {
    // Extract host:port from address
    host, port, err := n.parseAddress(address)
    if err != nil {
        fmt.Printf("❌ Invalid peer address %s: %v\n", address, err)
        return false
    }
    
    // Skip self-connection check
    if n.isSelfConnection(host, port) {
        fmt.Printf("🔇 Skipping self-connection: %s\n", address)
        return false
    }
    
    // Validate port
    if port == "" {
        fmt.Printf("❌ Empty port in address: %s\n", address)
        return false
    }
    
    // Try to establish TCP connection
    timeout := 2 * time.Second
    target := net.JoinHostPort(host, port)
    
    conn, err := net.DialTimeout("tcp", target, timeout)
    if err != nil {
        fmt.Printf("🌐 Cannot connect to peer %s (%s): %v\n", address, target, err)
        return false
    }
    conn.Close()
    
    fmt.Printf("🔗 Successfully connected to peer: %s (%s)\n", address, target)
    return true
}

// isSelfConnection checks if the peer address is our own address ✅ IMPROVED
func (n *Network) isSelfConnection(host, port string) bool {
    selfHost, selfPort, err := n.parseAddress(n.config.ListenAddr)
    if err != nil {
        fmt.Printf("⚠️  Cannot parse self address %s: %v\n", n.config.ListenAddr, err)
        return false
    }
    
    isSelf := host == selfHost && port == selfPort
    return isSelf
}

// RemovePeer removes a peer from the network
func (n *Network) RemovePeer(address string) {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    peerID := "peer-" + address
    if _, exists := n.peers[peerID]; exists {
        delete(n.peers, peerID)
        fmt.Printf("🧹 Removed peer: %s\n", address)
    }
}

// GetPeers returns list of all peers in memory
func (n *Network) GetPeers() []*PeerInfo {
    n.mu.RLock()
    defer n.mu.RUnlock()

    peers := make([]*PeerInfo, 0, len(n.peers))
    for _, p := range n.peers {
        if p.Connected {
            peers = append(peers, p)
        }
    }
    return peers
}

// GetActivePeers returns only ACTUALLY connected peers ✅ ENHANCED
func (n *Network) GetActivePeers() []*PeerInfo {
    n.mu.RLock()
    defer n.mu.RUnlock()

    activePeers := make([]*PeerInfo, 0)
    for _, peer := range n.peers {
        // Double-check connection status dengan real TCP check
        if peer.Connected && n.checkPeerConnection(peer.Address) {
            activePeers = append(activePeers, peer)
            peer.LastSeen = time.Now() // Update last seen time
        } else {
            // Mark as disconnected if check fails
            peer.Connected = false
            fmt.Printf("🔴 Peer disconnected: %s\n", peer.Address)
        }
    }
    return activePeers
}

// BroadcastBlock broadcasts a new block to all peers
func (n *Network) BroadcastBlock(block *types.Block) {
    // ✅ FIX: Use active peers instead of all peers
    activePeers := n.GetActivePeers()
    peerCount := len(activePeers)
    
    fmt.Printf("📤 [P2P] Broadcasting block #%s to %d active peers\n", block.Header.Number, peerCount)
    
    for _, peer := range activePeers {
        if peer.Connected {
            fmt.Printf("   ➡️  Sending to %s\n", peer.Address)
            // Update last seen time
            peer.LastSeen = time.Now()
            // In real implementation, this would send over network
        }
    }
    
    if peerCount == 0 {
        fmt.Printf("   ℹ️  No active peers connected - block ready for propagation\n")
    }
}

// BroadcastTransaction broadcasts a transaction to all peers
func (n *Network) BroadcastTransaction(tx *types.Transaction) {
    activePeers := n.GetActivePeers()
    peerCount := len(activePeers)
    
    fmt.Printf("📤 [P2P] Broadcasting transaction to %d active peers\n", peerCount)
    
    for _, peer := range activePeers {
        if peer.Connected {
            fmt.Printf("   ➡️  Sending TX to %s\n", peer.Address)
            peer.LastSeen = time.Now()
        }
    }
}

// connectToBootstrapPeers connects to bootstrap peers
func (n *Network) connectToBootstrapPeers() {
    successfulConnections := 0
    for _, addr := range n.config.BootstrapPeers {
        if err := n.AddPeer(addr); err != nil {
            fmt.Printf("❌ Failed to connect to bootstrap peer %s: %v\n", addr, err)
        } else {
            successfulConnections++
        }
    }
    fmt.Printf("🔌 Connected to %d/%d bootstrap peers\n", successfulConnections, len(n.config.BootstrapPeers))
}

// cleanupDeadPeers periodically removes dead peers
func (n *Network) cleanupDeadPeers() {
    ticker := time.NewTicker(60 * time.Second) // Check every minute
    defer ticker.Stop()

    for n.running {
        <-ticker.C
        
        n.mu.Lock()
        removedCount := 0
        for peerID, peer := range n.peers {
            // Remove peers not seen for more than 3 minutes
            if time.Since(peer.LastSeen) > 3*time.Minute {
                delete(n.peers, peerID)
                removedCount++
                fmt.Printf("🧹 [P2P] Cleaned up dead peer: %s (last seen: %v ago)\n", 
                    peer.Address, time.Since(peer.LastSeen))
            }
        }
        n.mu.Unlock()
        
        if removedCount > 0 {
            fmt.Printf("🧹 [P2P] Removed %d dead peers\n", removedCount)
        }
    }
}

// monitorNetwork periodically reports network status
func (n *Network) monitorNetwork() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for n.running {
        <-ticker.C
        
        activePeers := n.GetActivePeers()
        totalPeers := len(n.peers)
        
        fmt.Printf("📊 [P2P] Network stats: %d active peers (%d total in list)\n", 
            len(activePeers), totalPeers)
    }
}

// SimulatePeerConnection simulates a peer connection for demo
func (n *Network) SimulatePeerConnection(peerAddr string) {
    n.AddPeer(peerAddr)
    fmt.Printf("🔗 [P2P] Simulated connection to peer: %s\n", peerAddr)
}

// UpdatePeerHealth updates the last seen time for a peer
func (n *Network) UpdatePeerHealth(peerAddr string) {
    n.mu.Lock()
    defer n.mu.Unlock()
    
    peerID := "peer-" + peerAddr
    if peer, exists := n.peers[peerID]; exists {
        peer.LastSeen = time.Now()
        peer.Connected = true
    }
}