package main

import (
    "flag"
    "fmt"
    "math/big"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/selsichain/selsichain-core/core/blockchain"
    "github.com/selsichain/selsichain-core/core/consensus/hybrid"
    "github.com/selsichain/selsichain-core/core/types"
    "github.com/selsichain/selsichain-core/p2p/network"
)

func main() {
    // Cloud environment detection
    if os.Getenv("RAILWAY_STATIC_URL") != "" {
        fmt.Println("☁️  =================================")
        fmt.Println("☁️  RUNNING IN RAILWAY CLOUD")
        fmt.Println("☁️  =================================")
        fmt.Printf("☁️  Railway URL: %s\n", os.Getenv("RAILWAY_STATIC_URL"))
        fmt.Printf("☁️  Railway Environment: %s\n", os.Getenv("RAILWAY_ENVIRONMENT"))
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
    consensusConfig := &hybrid.Config{
        PowBlockInterval: 5,
        MiningDifficulty: big.NewInt(1000000),
        MinimumStake:     new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
        BlockTime:        12 * time.Second,
        RewardDistribution: hybrid.RewardConfig{
            MinerPercent:     45,
            StakerPercent:    45,
            EcosystemPercent: 7,
            BurnPercent:      3,
        },
    }

    // Initialize P2P network config
    networkConfig := &network.Config{
        ListenAddr: "/ip4/0.0.0.0/tcp/" + p2pPort,
        BootstrapPeers: []string{
            "/ip4/127.0.0.1/tcp/7691",
            "/ip4/127.0.0.1/tcp/7692",
        },
        ChainID:    769,  // Mainnet chain ID
        ProtocolID: "/selsichain",
    }

    // Testnet configuration
    if testnet {
        fmt.Println("🌐 TESTNET MODE ACTIVATED!")
        fmt.Println("🔧 Testnet Chain ID: 1337")
        fmt.Println("🎯 Testnet Validators: 5")
        
        // Adjust config for testnet
        consensusConfig.PowBlockInterval = 3  // More frequent PoW checkpoints
        consensusConfig.BlockTime = 10 * time.Second  // Faster blocks
        
        // Testnet-specific bootstrap peers
        networkConfig.BootstrapPeers = []string{
            "/ip4/127.0.0.1/tcp/7690",
            "/ip4/127.0.0.1/tcp/7691", 
            "/ip4/127.0.0.1/tcp/7692",
            "/ip4/127.0.0.1/tcp/7693",
            "/ip4/127.0.0.1/tcp/7694",
        }
        networkConfig.ChainID = 1337  // Testnet chain ID
    }

    consensusEngine := hybrid.NewHybridEngine(consensusConfig)

    // Initialize blockchain
    fmt.Println("🔄 Creating blockchain...")
    chain, err := blockchain.NewBlockchain(&blockchain.Config{
        DataDir: "./data",
    }, consensusEngine)

    if err != nil {
        fmt.Printf("❌ Error: %v\n", err)
        return
    }

    // Initialize P2P network
    fmt.Println("🌐 Initializing P2P Network...")
    p2pNetwork, err := network.NewNetwork(networkConfig, chain)
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
    activePeers := p2pNetwork.GetActivePeers()
    allPeers := p2pNetwork.GetPeers()
    
    fmt.Println("")
    fmt.Println("🎉 SELSIHAIN FULL NODE STARTED!")
    fmt.Println("===============================")
    fmt.Printf("⛓️  Current Block: #%s\n", currentBlock.Header.Number)
    fmt.Printf("📊 Total Blocks: %d\n", chain.GetBlockCount())
    fmt.Printf("🌐 P2P Address: %s\n", p2pNetwork.GetListenAddr())
    fmt.Printf("🆔 Peer ID: %s\n", p2pNetwork.GetPeerID())
    
    // REAL Peer Count vs Total Attempted
    fmt.Printf("👥 Connected Peers: %d/%d (active/total)\n", len(activePeers), len(allPeers))
    
    // Show peer connection details
    if len(activePeers) == 0 {
        fmt.Printf("   🔍 No active peer connections\n")
    } else {
        for _, peer := range activePeers {
            fmt.Printf("   ✅ %s\n", peer.Address)
        }
    }
    
    fmt.Println("===============================")
    fmt.Println("")
    fmt.Println("💡 CLI wallet not available in cloud deployment")
    fmt.Println("")

    // Demo: Create blocks
    fmt.Println("🧪 Creating demo blocks...")
    createDemoBlocks(chain, consensusEngine, p2pNetwork)
    
    fmt.Println("")
    fmt.Println("✅ Node is running and ready!")
    fmt.Println("⏳ Press Ctrl+C to shutdown")

    waitForShutdown(chain, p2pNetwork)
}

func createDemoBlocks(chain *blockchain.Blockchain, consensus *hybrid.HybridEngine, p2pNetwork *network.Network) {
    blockCount := 1
    
    for {
        fmt.Printf("\n🎯 Creating block #%d...\n", blockCount)
        
        currentBlock := chain.GetCurrentBlock()
        
        // Create sample transaction
        toAddr := types.Address{2}
        tx := &types.Transaction{
            Nonce:    uint64(blockCount),
            GasPrice: big.NewInt(1000000000),
            Gas:      21000,
            To:       &toAddr,
            Value:    new(big.Int).Mul(big.NewInt(10), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
            Data:     []byte{},
            Type:     types.TxRegular,
        }
        
        // Create block
        newBlock, err := consensus.CreateBlock(
            currentBlock,
            []*types.Transaction{tx},
            types.Address{byte(blockCount % 3 + 1)}, // Rotate between validators
        )
        
        if err == nil {
            // Verify and add block
            state := chain.GetStateDB()
            if consensus.VerifyBlock(newBlock, state) == nil {
                chain.AddBlock(newBlock)
                
                // FIX: Broadcast only to active peers
                p2pNetwork.BroadcastBlock(newBlock)
                fmt.Printf("✅ Block #%d created and broadcasted\n", blockCount)
                
                // Show milestone every 10 blocks
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

func waitForShutdown(chain *blockchain.Blockchain, p2pNetwork *network.Network) {
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    
    <-sigCh
    fmt.Println("")
    fmt.Println("🛑 Shutdown signal received...")
    
    p2pNetwork.Stop()
    chain.Close()
    
    fmt.Println("👋 SelsiChain node stopped gracefully")
}
