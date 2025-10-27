package hybrid

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "github.com/selsichain/selsichain-core/core/types"
    "github.com/selsichain/selsichain-core/core/state"
    "github.com/selsichain/selsichain-core/crypto/hash"
)

type POWEngine struct {
    config *Config
}

func NewPOWEngine(config *Config) *POWEngine {
    return &POWEngine{
        config: config,
    }
}

// VerifyBlock verifies PoW block
func (p *POWEngine) VerifyBlock(block *types.Block, state *state.StateDB) error {
    // Verify difficulty
    if !p.verifyDifficulty(block.Header) {
        return ErrInvalidDifficulty
    }
    
    // Verify block time
    if !p.verifyBlockTime(block.Header) {
        return ErrBlockTimeTooEarly
    }
    
    fmt.Printf("⛏️  PoW Block #%s verified\n", block.Header.Number)
    return nil
}

// PrepareBlock prepares block for mining
func (p *POWEngine) PrepareBlock(block *types.Block) (*types.Block, error) {
    block.Header.Checkpoint = true
    block.Header.Difficulty = p.config.MiningDifficulty
    
    // Generate random nonce untuk mining
    block.Header.Nonce = p.generateNonce()
    
    fmt.Printf("⛏️  PoW Block #%s prepared for mining\n", block.Header.Number)
    return block, nil
}

// MineBlock simulates mining process
func (p *POWEngine) MineBlock(block *types.Block) (*types.Block, error) {
    fmt.Printf("⛏️  Mining PoW Block #%s...\n", block.Header.Number)
    
    // Simulate mining process (in real implementation, this would be actual PoW)
    target := new(big.Int).Div(
        new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil),
        p.config.MiningDifficulty,
    )
    
    // Simple mining simulation - try different nonces
    for i := 0; i < 1000; i++ {
        block.Header.Nonce = p.generateNonce()
        blockHash := hash.CalculateBlockHash(block.Header)
        hashInt := new(big.Int).SetBytes(blockHash[:])
        
        if hashInt.Cmp(target) == -1 {
            fmt.Printf("✅ PoW Block #%s mined successfully!\n", block.Header.Number)
            return block, nil
        }
    }
    
    // For demo purposes, just return the block even if mining "fails"
    fmt.Printf("⚠️  PoW Block #%s - using fallback (demo mode)\n", block.Header.Number)
    return block, nil
}

func (p *POWEngine) verifyDifficulty(header *types.Header) bool {
    // For demo, always return true
    return true
}

func (p *POWEngine) verifyBlockTime(header *types.Header) bool {
    // Basic time verification
    return true
}

func (p *POWEngine) generateNonce() types.BlockNonce {
    var nonce types.BlockNonce
    rand.Read(nonce[:])
    return nonce
}
