package wallet

import (
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "math/big"

    "github.com/selsichain/selsichain-core/core/blockchain"
    "github.com/selsichain/selsichain-core/core/types"
)

// SimpleWallet adalah simplified wallet untuk demo
type SimpleWallet struct {
    chain      *blockchain.Blockchain
    accounts   map[string]*SimpleAccount
    defaultAcc string
}

// SimpleAccount represents a simple account
type SimpleAccount struct {
    Address    types.Address
    PrivateKey string
    Balance    *big.Int
    Stake      *big.Int
}

// NewSimpleWallet creates a new simple wallet
func NewSimpleWallet(chain *blockchain.Blockchain) *SimpleWallet {
    return &SimpleWallet{
        chain:    chain,
        accounts: make(map[string]*SimpleAccount),
    }
}

// CreateAccount creates a new account
func (sw *SimpleWallet) CreateAccount() (*SimpleAccount, error) {
    // Generate random private key (simplified)
    privateKey := make([]byte, 32)
    rand.Read(privateKey)
    privateKeyHex := hex.EncodeToString(privateKey)

    // Generate address from private key (simplified)
    address := sw.generateAddress(privateKeyHex)

    account := &SimpleAccount{
        Address:    address,
        PrivateKey: privateKeyHex,
        Balance:    big.NewInt(0),
        Stake:      big.NewInt(0),
    }

    // Add some initial balance for demo
    account.Balance = new(big.Int).Mul(big.NewInt(1000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

    sw.accounts[account.GetAddressHex()] = account

    // Set as default if first account
    if sw.defaultAcc == "" {
        sw.defaultAcc = account.GetAddressHex()
    }

    fmt.Printf("✅ Account created: %s\n", account.GetAddressHex())
    fmt.Printf("   Private Key: %s (SAVE THIS SECURELY!)\n", privateKeyHex)
    fmt.Printf("   Balance: 1000 SELSI (demo)\n")

    return account, nil
}

// ImportAccount imports an account from private key
func (sw *SimpleWallet) ImportAccount(privateKeyHex string) (*SimpleAccount, error) {
    address := sw.generateAddress(privateKeyHex)

    account := &SimpleAccount{
        Address:    address,
        PrivateKey: privateKeyHex,
        Balance:    big.NewInt(0),
        Stake:      big.NewInt(0),
    }

    // Add some initial balance for demo
    account.Balance = new(big.Int).Mul(big.NewInt(500), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

    sw.accounts[account.GetAddressHex()] = account

    // Set as default if no default
    if sw.defaultAcc == "" {
        sw.defaultAcc = account.GetAddressHex()
    }

    fmt.Printf("✅ Account imported: %s\n", account.GetAddressHex())
    fmt.Printf("   Balance: 500 SELSI (demo)\n")

    return account, nil
}

// GetBalance returns the balance of an account
func (sw *SimpleWallet) GetBalance(address types.Address) (*big.Int, error) {
    // For demo, use the wallet balance
    // In real implementation, this would query the blockchain
    account := sw.accounts[hex.EncodeToString(address[:])]
    if account == nil {
        return big.NewInt(0), nil
    }
    return account.Balance, nil
}

// GetStake returns the staked amount
func (sw *SimpleWallet) GetStake(address types.Address) (*big.Int, error) {
    account := sw.accounts[hex.EncodeToString(address[:])]
    if account == nil {
        return big.NewInt(0), nil
    }
    return account.Stake, nil
}

// SendTransaction sends SELSI tokens
func (sw *SimpleWallet) SendTransaction(from, to types.Address, amount *big.Int) error {
    fromAccount := sw.accounts[hex.EncodeToString(from[:])]
    if fromAccount == nil {
        return fmt.Errorf("sender account not found")
    }

    toAccount := sw.accounts[hex.EncodeToString(to[:])]
    if toAccount == nil {
        // Create recipient account if doesn't exist
        toAccount = &SimpleAccount{
            Address: to,
            Balance: big.NewInt(0),
            Stake:   big.NewInt(0),
        }
        sw.accounts[hex.EncodeToString(to[:])] = toAccount
    }

    // Check balance
    if fromAccount.Balance.Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance")
    }

    // Transfer
    fromAccount.Balance.Sub(fromAccount.Balance, amount)
    toAccount.Balance.Add(toAccount.Balance, amount)

    fmt.Printf("✅ Transfer successful!\n")
    fmt.Printf("   From: %s\n", hex.EncodeToString(from[:4]))
    fmt.Printf("   To: %s\n", hex.EncodeToString(to[:4]))
    fmt.Printf("   Amount: %s SELSI\n", new(big.Int).Div(amount, big.NewInt(1e18)))

    return nil
}

// StakeTokens stakes SELSI tokens
func (sw *SimpleWallet) StakeTokens(from types.Address, amount *big.Int) error {
    account := sw.accounts[hex.EncodeToString(from[:])]
    if account == nil {
        return fmt.Errorf("account not found")
    }

    // Check balance
    if account.Balance.Cmp(amount) < 0 {
        return fmt.Errorf("insufficient balance")
    }

    // Stake
    account.Balance.Sub(account.Balance, amount)
    account.Stake.Add(account.Stake, amount)

    fmt.Printf("✅ Staking successful!\n")
    fmt.Printf("   Account: %s\n", hex.EncodeToString(from[:4]))
    fmt.Printf("   Staked: %s SELSI\n", new(big.Int).Div(amount, big.NewInt(1e18)))
    fmt.Printf("   New Stake: %s SELSI\n", new(big.Int).Div(account.Stake, big.NewInt(1e18)))

    return nil
}

// ListAccounts lists all accounts
func (sw *SimpleWallet) ListAccounts() []*SimpleAccount {
    accounts := make([]*SimpleAccount, 0, len(sw.accounts))
    for _, account := range sw.accounts {
        accounts = append(accounts, account)
    }
    return accounts
}

// GetDefaultAccount returns the default account
func (sw *SimpleWallet) GetDefaultAccount() *SimpleAccount {
    if sw.defaultAcc == "" {
        return nil
    }
    return sw.accounts[sw.defaultAcc]
}

// SetDefaultAccount sets the default account
func (sw *SimpleWallet) SetDefaultAccount(address types.Address) error {
    addressHex := hex.EncodeToString(address[:])
    if sw.accounts[addressHex] == nil {
        return fmt.Errorf("account not found")
    }
    sw.defaultAcc = addressHex
    return nil
}

// generateAddress generates a simple address from private key
func (sw *SimpleWallet) generateAddress(privateKeyHex string) types.Address {
    // Simple address generation for demo
    // In real implementation, use proper cryptographic address generation
    hash := hex.EncodeToString([]byte(privateKeyHex))
    var address types.Address
    copy(address[:], hash[:20])
    return address
}

// GetAddressHex returns address as hex string
func (sa *SimpleAccount) GetAddressHex() string {
    return hex.EncodeToString(sa.Address[:])
}

// GetAccountInfo returns account information
func (sw *SimpleWallet) GetAccountInfo(address types.Address) *AccountInfo {
    account := sw.accounts[hex.EncodeToString(address[:])]
    if account == nil {
        return nil
    }

    return &AccountInfo{
        Address: address,
        Balance: account.Balance,
        Stake:   account.Stake,
        Nonce:   0, // Simplified
    }
}

// AccountInfo represents account information
type AccountInfo struct {
    Address types.Address
    Balance *big.Int
    Stake   *big.Int
    Nonce   uint64
}
