package types

import "math/big"

// Hash represents 32-byte hash
type Hash [32]byte

// Address represents 20-byte address  
type Address [20]byte

// BlockNonce represents 8-byte nonce
type BlockNonce [8]byte

// Block represents a SelsiChain block
type Block struct {
    Header       *Header
    Transactions []*Transaction
    Votes        []*Vote  // Untuk PoS
}

// Header represents block header
type Header struct {
    ParentHash   Hash
    Coinbase     Address  // Miner/validator address
    Root         Hash
    TxHash       Hash
    Difficulty   *big.Int
    Number       *big.Int
    Time         uint64
    Extra        []byte
    MixDigest    Hash
    Nonce        BlockNonce
    
    // SelsiChain Hybrid Fields
    Validator    Address  // PoS validator
    StakeHash    Hash     // Stake merkle root  
    Checkpoint   bool     // Is PoW checkpoint block?
}

// Transaction represents a transaction
type Transaction struct {
    Nonce    uint64
    GasPrice *big.Int
    Gas      uint64
    To       *Address    // Pointer karena bisa nil (contract creation)
    Value    *big.Int
    Data     []byte      // Renamed from Input untuk consistency
    V, R, S  *big.Int
    
    // SelsiChain specific
    Type     TxType
}

// Vote represents a PoS vote
type Vote struct {
    Validator   Address
    BlockHash   Hash
    Decision    bool
    Signature   []byte
    Timestamp   int64
}

// TxType represents transaction type
type TxType uint8

const (
    TxRegular TxType = iota
    TxStaking
    TxUnstaking
    TxVoting
)
