package blockchain

import (
    "fmt"
    "math/big"
    "time"
    "github.com/selsichain/selsichain-core/core/types"
    "github.com/selsichain/selsichain-core/core/state"
    "github.com/selsichain/selsichain-core/core/consensus/hybrid"
)

type Blockchain struct {
    genesis   *types.Block
    current   *types.Block
    blocks    map[string]*types.Block
    state     *state.StateDB
    consensus *hybrid.HybridEngine
    config    *Config
}

type Config struct {
    DataDir string
}

func NewBlockchain(config *Config, consensus *hybrid.HybridEngine) (*Blockchain, error) {
    bc := &Blockchain{
        blocks:    make(map[string]*types.Block),
        state:     state.NewStateDB(),
        consensus: consensus,
        config:    config,
    }
    
    if err := bc.initGenesis(); err != nil {
        return nil, err
    }
    
    return bc, nil
}

func (bc *Blockchain) initGenesis() error {
    bc.genesis = &types.Block{
        Header: &types.Header{
            ParentHash: types.Hash{},
            Number:     big.NewInt(0),
            Time:       uint64(time.Now().Unix()),
            Difficulty: big.NewInt(1000),
        },
        Transactions: []*types.Transaction{},
    }
    
    bc.blocks[bc.CalculateHash(bc.genesis.Header)] = bc.genesis
    bc.current = bc.genesis
    
    // Initialize genesis accounts with balances and stakes
    genesisAddr1 := types.Address{1}
    genesisAddr2 := types.Address{2} 
    genesisAddr3 := types.Address{3}
    
    // Set initial balances
    initialBalance := new(big.Int).Mul(big.NewInt(1000000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
    bc.state.SetBalance(genesisAddr1, initialBalance)
    bc.state.SetBalance(genesisAddr2, initialBalance)
    bc.state.SetBalance(genesisAddr3, initialBalance)
    
    // Set initial stakes for PoS validators
    stake1 := new(big.Int).Mul(big.NewInt(5000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
    stake2 := new(big.Int).Mul(big.NewInt(3000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
    stake3 := new(big.Int).Mul(big.NewInt(7000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
    
    bc.state.SetStake(genesisAddr1, stake1)
    bc.state.SetStake(genesisAddr2, stake2) 
    bc.state.SetStake(genesisAddr3, stake3)
    
    fmt.Println("✅ Genesis block created with 3 validator accounts")
    fmt.Printf("   Validator 1: %x - Stake: %s SELSI\n", genesisAddr1[:4], stake1)
    fmt.Printf("   Validator 2: %x - Stake: %s SELSI\n", genesisAddr2[:4], stake2)
    fmt.Printf("   Validator 3: %x - Stake: %s SELSI\n", genesisAddr3[:4], stake3)
    
    return nil
}

func (bc *Blockchain) CalculateHash(header *types.Header) string {
    return fmt.Sprintf("block-%d-%d", header.Number, header.Time)
}

func (bc *Blockchain) AddBlock(block *types.Block) error {
    hash := bc.CalculateHash(block.Header)
    bc.blocks[hash] = block
    bc.current = block
    
    // Apply rewards from consensus
    rewards := bc.consensus.CalculateRewards(block, bc.state)
    for addr, reward := range rewards {
        bc.state.AddBalance(addr, reward)
        fmt.Printf("💰 Rewarded %x: +%s SELSI\n", addr[:4], reward)
    }
    
    fmt.Printf("✅ Block #%s added to chain\n", block.Header.Number)
    return nil
}

func (bc *Blockchain) GetCurrentBlock() *types.Block {
    return bc.current
}

func (bc *Blockchain) Close() {
    fmt.Println("📦 Blockchain closed")
}

func (bc *Blockchain) GetBlockCount() int {
    return len(bc.blocks)
}

func (bc *Blockchain) GetStateDB() *state.StateDB {
    return bc.state
}
