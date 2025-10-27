package state

import (
    "math/big"
    "github.com/selsichain/selsichain-core/core/types"
)

// StateDB manages the state of accounts and stakes
type StateDB struct {
    accounts map[types.Address]*Account
    stakes   map[types.Address]*big.Int
}

// Account represents a user account
type Account struct {
    Balance *big.Int
    Nonce   uint64
    Code    []byte
}

// NewStateDB creates a new state database
func NewStateDB() *StateDB {
    return &StateDB{
        accounts: make(map[types.Address]*Account),
        stakes:   make(map[types.Address]*big.Int),
    }
}

// GetBalance returns the balance of an address
func (s *StateDB) GetBalance(address types.Address) *big.Int {
    if account, exists := s.accounts[address]; exists {
        return new(big.Int).Set(account.Balance)
    }
    return big.NewInt(0)
}

// SetBalance sets the balance of an address
func (s *StateDB) SetBalance(address types.Address, amount *big.Int) {
    if _, exists := s.accounts[address]; !exists {
        s.accounts[address] = &Account{}
    }
    s.accounts[address].Balance = new(big.Int).Set(amount)
}

// GetStake returns the staked amount of an address
func (s *StateDB) GetStake(address types.Address) *big.Int {
    if stake, exists := s.stakes[address]; exists {
        return new(big.Int).Set(stake)
    }
    return big.NewInt(0)
}

// SetStake sets the staked amount of an address
func (s *StateDB) SetStake(address types.Address, amount *big.Int) {
    s.stakes[address] = new(big.Int).Set(amount)
}

// GetNonce returns the nonce of an address
func (s *StateDB) GetNonce(address types.Address) uint64 {
    if account, exists := s.accounts[address]; exists {
        return account.Nonce
    }
    return 0
}

// SetNonce sets the nonce of an address
func (s *StateDB) SetNonce(address types.Address, nonce uint64) {
    if _, exists := s.accounts[address]; !exists {
        s.accounts[address] = &Account{}
    }
    s.accounts[address].Nonce = nonce
}

// AddBalance adds amount to the balance of an address
func (s *StateDB) AddBalance(address types.Address, amount *big.Int) {
    current := s.GetBalance(address)
    s.SetBalance(address, new(big.Int).Add(current, amount))
}

// SubBalance subtracts amount from the balance of an address
func (s *StateDB) SubBalance(address types.Address, amount *big.Int) {
    current := s.GetBalance(address)
    s.SetBalance(address, new(big.Int).Sub(current, amount))
}

// Exist checks if an address exists
func (s *StateDB) Exist(address types.Address) bool {
    _, exists := s.accounts[address]
    return exists
}

// Empty checks if an address is empty
func (s *StateDB) Empty(address types.Address) bool {
    if !s.Exist(address) {
        return true
    }
    account := s.accounts[address]
    return account.Balance.Cmp(big.NewInt(0)) == 0 &&
        account.Nonce == 0 &&
        len(account.Code) == 0
}
