package main

import (
    "flag"
    "fmt"
    "math/big"
    "os"
    "os/signal"
    "syscall"
    "time"
)

// Config untuk blockchain
type BlockchainConfig struct {
    DataDir string
}

// Config untuk hybrid consensus
type HybridConfig struct {
    PowBlockInterval int
    MiningDifficulty *big.Int
    MinimumStake     *big.Int
    BlockTime        time.Duration
    RewardDistribution RewardConfig
}

type RewardConfig struct {
    MinerPercent     int
    StakerPercent    int
    EcosystemPercent int
    BurnPercent      int
}

// Config untuk network
type NetworkConfig struct {
    ListenAddr     string
    BootstrapPeers []string
    ChainID        uint64
    ProtocolID     string
}

// Mock types untuk demo
type Address [20]byte
type Block struct {
    Header BlockHeader
}
type BlockHeader struct {
    Number *big.Int
}
type Transaction struct {
    Nonce    uint64
    GasPrice *big.Int
    Gas      uint64
    To       *Address
    Value    *big.Int
    Data     []byte
    Type     string
}

// Mock functions untuk demo
func NewBlockchain(config *BlockchainConfig, engine interface{}) (*Blockchain, error) {
    fmt.Println("📦 Creating new blockchain...")
    return &Blockchain{currentBlock: &Block{Header: BlockHeader{Number: big.NewInt(0)}}}, nil
}

type Blockchain struct {
    currentBlock *Block
}

func (bc *Blockchain) GetCurrentBlock() *Block {
    return bc.currentBlock
}

func (bc *Blockchain) GetBlockCount() int {
    return 1
}

func (bc *Blockchain) AddBlock(block *Block) {
    bc.currentBlock = block
    fmt.Printf("✅ Added block #%s to chain\n", block.Header.Number)
}

func (bc *Blockchain) Close() {
    fmt.Println("💾 Closing blockchain...")
}

func (bc *Blockchain) GetStateDB() interface{} {
    return nil
}

func NewHybridEngine(config *HybridConfig) *HybridEngine {
    fmt.Println("⚙️  Creating hybrid consensus engine...")
    return &HybridEngine{config: config}
}

type HybridEngine struct {
    config *HybridConfig
}

func (he *HybridEngine) CreateBlock(parent *Block, txs []*Transaction, validator Address) (*Block, error) {
    newBlock := &Block{
        Header: BlockHeader{
            Number: new(big.Int).Add(parent.Header.Number, big.NewInt(1)),
        },
    }
    fmt.Printf("🎯 Created block #%s\n", newBlock.Header.Number)
    return newBlock, nil
}

func (he *HybridEngine) VerifyBlock(block *Block, state interface{}) error {
    fmt.Printf("🔍 Verifying block #%s... ✅\n", block.Header.Number)
    return nil
}

func NewNetwork(config *NetworkConfig, chain *Blockchain) (*Network, error) {
    fmt.Println("🌐 Creating P2P network...")
    return &Network{config: config}, nil
}

type Network struct {
    config *NetworkConfig
}

func (n *Network) Start() error {
    fmt.Printf("📍 Listening on: %s\n", n.config.ListenAddr)
    fmt.Printf("🔗 Protocol: %s\n", n.config.ProtocolID)
    fmt.Printf("⛓️  Chain ID: %d\n", n.config.ChainID)
    fmt.Printf("🔌 Connecting to %d bootstrap peers...\n", len(n.config.BootstrapPeers))
    
    for _, peer := range n.config.BootstrapPeers {
        fmt.Printf("   🔗 Attempting: %s\n", peer)
    }
    
    fmt.Println("✅ P2P Network started successfully!")
    return nil
}

func (n *Network) Stop() {
    fmt.Println("🛑 Stopping P2P network...")
}

func (n *Network) GetListenAddr() string {
    return n.config.ListenAddr
}

func (n *Network) GetPeerID() string {
    return "selsichain-node-" + n.config.ListenAddr
}

func (n *Network) GetActivePeers() []interface{} {
    return []interface{}{}
}

func (n *Network) GetPeers() []interface{} {
    return []interface{}{}
}

func (n *Network) BroadcastBlock(block *Block) {
    fmt.Printf("📤 Broadcasting block #%s to peers...\n", block.Header.Number)
}

func main() {
    // Cloud environment detection
    if os.Getenv("RAILWAY_STATIC_URL") != "" {
        fmt.Println("☁️  =================================")
        fmt.Println("☁️  RUNNING IN RAILWAY CLOUD")
        fmt.Println("☁️  =================================")
        fmt.Printf("☁️  Railway URL: %s\n", os.Getenv("RAILWAY_STATIC_URL"))
        if os.Getenv("RAILWAY_GIT_COMMIT_SHA") != "" {
            fmt.Printf("☁️  Deployment SHA: %s\n", os.Getenv("RAILWAY_GIT_COMMIT_SHA"))
        }
        fmt.Println("☁️  =================================")
    }

    // Parse command line flags
    p2pPort := flag.String("p2p-port", "7690", "P2P network port")
    testnet := flag.Bool("testnet", false, "Enable testnet mode")
    flag.Parse()

    startFullNode(*p2pPort, *testnet)
}

func startFullNode(p2pPort string, testnet bool) {
    // Use PORT from environment if running in cloud
    if envPort := os.Getenv("PORT"); envPort != "" && p2pPort == "7690" {
        p2pPort = envPort
        fmt.Printf("☁️  Using cloud PORT: %s\n", p2pPort)
    }

    fmt.Println("")
    fmt.Println("🚀 Starting SelsiChain Full Node...")
    fmt.Println("")

    // Initialize hybrid consensus
    fmt.Println("🔄 Initializing Hybrid Consensus Engine...")
    consensusConfig := &HybridConfig{
        PowBlockInterval: 5,
        MiningDifficulty: big.NewInt(1000000),
        MinimumStake:     new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
        BlockTime:        12 * time.Second,
        RewardDistribution: RewardConfig{
            MinerPercent:     45,
            StakerPercent:    45,
            EcosystemPercent: 7,
            BurnPercent:      3,
        },
    }

    // Initialize P2P network config
    networkConfig := &NetworkConfig{
        ListenAddr: "/ip4/0.0.0.0/tcp/" + p2pPort,
        BootstrapPeers: []string{
            "/ip4/127.0.0.1/tcp/7691",
            "/ip4/127.0.0.1/tcp/7692",
        },
        ChainID:    769,
        ProtocolID: "/selsichain",
    }

    // Testnet configuration
    if testnet {
        fmt.Println("🌐 TESTNET MODE ACTIVATED!")
        fmt.Println("🔧 Testnet Chain ID: 1337")
        fmt.Println("🎯 Testnet Validators: 5")
        
        consensusConfig.PowBlockInterval = 3
        consensusConfig.BlockTime = 10 * time.Second
        
        networkConfig.BootstrapPeers = []string{
            "/ip4/127.0.0.1/tcp/7690",
            "/ip4/127.0.0.1/tcp/7691", 
            "/ip4/127.0.0.1/tcp/7692",
            "/ip4/127.0.0.1/tcp/7693",
            "/ip4/127.0.0.1/tcp/7694",
        }
        networkConfig.ChainID = 1337
    }

    consensusEngine := NewHybridEngine(consensusConfig)

    // Initialize blockchain
    fmt.Println("🔄 Creating blockchain...")
    chain, err := NewBlockchain(&BlockchainConfig{
        DataDir: "./data",
    }, consensusEngine)

    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    // Initialize P2P network
    fmt.Println("🌐 Initializing P2P Network...")
    p2pNetwork, err := NewNetwork(networkConfig, chain)
    if err != nil {
        fmt.Printf("❌ Failed to initialize P2P network: %v\n", err)
        return
    }

    // Start P2P network
    if err := p2pNetwork.Start(); err != nil {
        fmt.Printf("❌ Failed to start P2P network: %v\n", err)
        return
    }

    // Display node info
    currentBlock := chain.GetCurrentBlock()
    
    fmt.Println("")
    fmt.Println("🎉 SELSIHAIN FULL NODE STARTED!")
    fmt.Println("===============================")
    fmt.Printf("⛓️  Current Block: #%s\n", currentBlock.Header.Number)
    fmt.Printf("📊 Total Blocks: %d\n", chain.GetBlockCount())
    fmt.Printf("🌐 P2P Address: %s\n", p2pNetwork.GetListenAddr())
    fmt.Printf("🆔 Peer ID: %s\n", p2pNetwork.GetPeerID())
    fmt.Printf("👥 Connected Peers: %d\n", len(p2pNetwork.GetActivePeers()))
    fmt.Println("===============================")
    fmt.Println("")
    fmt.Println("💡 Cloud deployment active - Demo mode")
    fmt.Println("")

    // Demo: Create blocks
    fmt.Println("🧪 Creating demo blocks...")
    createDemoBlocks(chain, consensusEngine, p2pNetwork)
    
    fmt.Println("")
    fmt.Println("✅ Node is running and ready!")
    fmt.Println("⏳ Press Ctrl+C to shutdown")

    waitForShutdown(chain, p2pNetwork)
}

func createDemoBlocks(chain *Blockchain, consensus *HybridEngine, p2pNetwork *Network) {
    blockCount := 1
    
    for {
        fmt.Printf("\n🎯 Creating block #%d...\n", blockCount)
        
        currentBlock := chain.GetCurrentBlock()
        
        // Create sample transaction
        toAddr := Address{2}
        tx := &Transaction{
            Nonce:    uint64(blockCount),
            GasPrice: big.NewInt(1000000000),
            Gas:      21000,
            To:       &toAddr,
            Value:    new(big.Int).Mul(big.NewInt(10), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
            Data:     []byte{},
            Type:     "regular",
        }
        
        // Create block
        newBlock, err := consensus.CreateBlock(
            currentBlock,
            []*Transaction{tx},
            Address{byte(blockCount % 3 + 1)},
        )
        
        if err == nil {
            if consensus.VerifyBlock(newBlock, nil) == nil {
                chain.AddBlock(newBlock)
                p2pNetwork.BroadcastBlock(newBlock)
                fmt.Printf("✅ Block #%d created and broadcasted\n", blockCount)
                
                if blockCount%10 == 0 {
                    fmt.Printf("🎉 Milestone: %d blocks produced!\n", blockCount)
                }
            }
        }
        
        // Wait 15 seconds before next block
        time.Sleep(15 * time.Second)
        blockCount++
    }
}

func waitForShutdown(chain *Blockchain, p2pNetwork *Network) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigCh
    fmt.Println("")
    fmt.Println("🛑 Shutdown signal received...")
    
    p2pNetwork.Stop()
    chain.Close()
    
    fmt.Println("👋 SelsiChain node stopped gracefully")
}
