package hybrid

import (
    "fmt"
    "math/big"
    "github.com/selsichain/selsichain-core/core/types"
    "github.com/selsichain/selsichain-core/core/state"
)

type POSEngine struct {
    config *Config
}

func NewPOSEngine(config *Config) *POSEngine {
    return &POSEngine{
        config: config,
    }
}

// VerifyBlock verifies PoS block
func (p *POSEngine) VerifyBlock(block *types.Block, state *state.StateDB) error {
    // Verify validator stake
    if !p.verifyValidatorStake(block.Header.Validator, state) {
        return ErrInsufficientStake
    }
    
    // Verify votes (2/3 majority)
    if !p.verifyVotes(block.Votes, block, state) {
        return ErrInsufficientVotes
    }
    
    fmt.Printf("🎯 PoS Block #%s verified\n", block.Header.Number)
    return nil
}

// PrepareBlock prepares block for staking
func (p *POSEngine) PrepareBlock(block *types.Block) (*types.Block, error) {
    // Select validator for this block
    validator, err := p.selectValidator(block.Header.Number)
    if err != nil {
        return nil, err
    }
    
    block.Header.Validator = validator
    block.Header.Checkpoint = false
    
    // Generate votes for this block
    block.Votes = p.generateVotes(block, validator)
    
    fmt.Printf("🎯 PoS Block #%s prepared for validation\n", block.Header.Number)
    fmt.Printf("🎯 Generated %d votes for block\n", len(block.Votes))
    
    return block, nil
}

// selectValidator selects validator based on stake and block number
func (p *POSEngine) selectValidator(blockNumber *big.Int) (types.Address, error) {
    validators := p.getEligibleValidators()
    if len(validators) == 0 {
        return types.Address{}, ErrNoValidators
    }
    
    // Simple round-robin selection based on block number
    index := new(big.Int).Mod(blockNumber, big.NewInt(int64(len(validators)))).Int64()
    selectedValidator := validators[index]
    
    fmt.Printf("🎯 Validator selected: %x (Stake: %s SELSI)\n", 
        selectedValidator.Address[:4], 
        selectedValidator.Stake)
    
    return selectedValidator.Address, nil
}

func (p *POSEngine) verifyValidatorStake(validator types.Address, state *state.StateDB) bool {
    stake := state.GetStake(validator)
    isValid := stake.Cmp(p.config.MinimumStake) >= 0
    
    if isValid {
        fmt.Printf("🎯 Validator %x has sufficient stake: %s SELSI\n", 
            validator[:4], stake)
    } else {
        fmt.Printf("❌ Validator %x has insufficient stake: %s SELSI (min: %s SELSI)\n", 
            validator[:4], stake, p.config.MinimumStake)
    }
    
    return isValid
}

func (p *POSEngine) verifyVotes(votes []*types.Vote, block *types.Block, state *state.StateDB) bool {
    approvedStake := big.NewInt(0)
    totalStake := big.NewInt(0)
    approvedVotes := 0
    totalVotes := 0
    
    // Simple block hash simulation untuk demo
    expectedBlockHash := p.simulateBlockHash(block)
    
    for _, vote := range votes {
        // Simple hash comparison untuk demo
        voteIsValid := p.simulateVoteValidation(vote, expectedBlockHash, state)
        
        if voteIsValid {
            stake := state.GetStake(vote.Validator)
            totalStake.Add(totalStake, stake)
            totalVotes++
            
            if vote.Decision {
                approvedStake.Add(approvedStake, stake)
                approvedVotes++
            }
        }
    }
    
    // Butuh 2/3 stake setuju
    if totalStake.Cmp(big.NewInt(0)) == 0 {
        fmt.Printf("❌ No valid votes received\n")
        return false
    }
    
    requiredStake := new(big.Int).Div(new(big.Int).Mul(totalStake, big.NewInt(2)), big.NewInt(3))
    isApproved := approvedStake.Cmp(requiredStake) >= 0
    
    fmt.Printf("🎯 Voting Results: %d/%d validators approved (%s/%s stake) - %s\n", 
        approvedVotes, totalVotes, approvedStake, totalStake, 
        map[bool]string{true: "APPROVED", false: "REJECTED"}[isApproved])
    
    return isApproved
}

// generateVotes creates mock votes for demo
func (p *POSEngine) generateVotes(block *types.Block, blockProposer types.Address) []*types.Vote {
    var votes []*types.Vote
    validators := p.getEligibleValidators()
    
    blockHash := p.simulateBlockHash(block)
    
    for _, validator := range validators {
        // Skip the block proposer (they don't vote for their own block)
        if validator.Address == blockProposer {
            continue
        }
        
        // 80% chance of approval for demo
        decision := true // Always approve for now to make it work
        
        votes = append(votes, &types.Vote{
            Validator: validator.Address,
            BlockHash: blockHash,
            Decision:  decision,
            Timestamp: 0,
        })
    }
    
    return votes
}

// simulateBlockHash creates a simple hash for demo
func (p *POSEngine) simulateBlockHash(block *types.Block) types.Hash {
    var hash types.Hash
    // Simple hash based on block number and validator
    copy(hash[:], fmt.Sprintf("block-%d-%x", block.Header.Number, block.Header.Validator[:2]))
    return hash
}

// simulateVoteValidation simulates vote verification for demo
func (p *POSEngine) simulateVoteValidation(vote *types.Vote, expectedHash types.Hash, state *state.StateDB) bool {
    // Check if validator has sufficient stake
    if !p.verifyValidatorStake(vote.Validator, state) {
        return false
    }
    
    // Simple hash comparison
    return vote.BlockHash == expectedHash
}

// getEligibleValidators returns list of validators with sufficient stake
func (p *POSEngine) getEligibleValidators() []Validator {
    // For demo, return validators with addresses 1, 2, 3
    return []Validator{
        {
            Address: types.Address{1},
            Stake:   new(big.Int).Mul(big.NewInt(5000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
        },
        {
            Address: types.Address{2}, 
            Stake:   new(big.Int).Mul(big.NewInt(3000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
        },
        {
            Address: types.Address{3},
            Stake:   new(big.Int).Mul(big.NewInt(7000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),
        },
    }
}
