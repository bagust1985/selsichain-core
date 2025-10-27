package hybrid

import (
    "math/big"
    "time"
    "github.com/selsichain/selsichain-core/core/types"
)

type Config struct {
    // PoW Configuration
    PowBlockInterval   uint64        // Setiap 100 block
    MiningDifficulty   *big.Int      // Difficulty target
    PowReward          *big.Int      // Reward untuk miner
    
    // PoS Configuration  
    MinimumStake       *big.Int      // Minimum stake required
    StakingPeriod      time.Duration // Lock period
    PosReward          *big.Int      // Reward untuk staker
    
    // Hybrid Configuration
    BlockTime          time.Duration // 12 detik
    RewardDistribution RewardConfig
}

type RewardConfig struct {
    MinerPercent     int // 45%
    StakerPercent    int // 45%
    EcosystemPercent int // 7%
    BurnPercent      int // 3%
}

type Validator struct {
    Address types.Address
    Stake   *big.Int
    Power   int64 // Voting power
}
