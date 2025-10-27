
"encoding/hex"

package main

import (
    "bufio"
    "fmt"
    "math/big"
    "os"
    "strconv"
    "strings"

    "github.com/selsichain/selsichain-core/core/blockchain"
    "github.com/selsichain/selsichain-core/core/consensus/hybrid"
    "github.com/selsichain/selsichain-core/crypto/wallet"
    "github.com/selsichain/selsichain-core/core/types"
)

// CLI represents the command-line interface
type CLI struct {
    wallet    *wallet.Wallet
    chain     *blockchain.Blockchain
    consensus *hybrid.HybridEngine
    running   bool
}

// NewCLI creates a new CLI
func NewCLI(chain *blockchain.Blockchain, consensus *hybrid.HybridEngine) *CLI {
    wallet := wallet.NewWallet("./keys", chain)
    
    return &CLI{
        wallet:    wallet,
        chain:     chain,
        consensus: consensus,
        running:   true,
    }
}

// Start starts the CLI interface
func (cli *CLI) Start() {
    fmt.Println("")
    fmt.Println("💰 SELSIHAIN WALLET CLI")
    fmt.Println("======================")
    fmt.Println("Type 'help' for available commands")
    fmt.Println("")

    scanner := bufio.NewScanner(os.Stdin)

    for cli.running {
        fmt.Print("selsichain> ")
        
        if !scanner.Scan() {
            break
        }

        input := strings.TrimSpace(scanner.Text())
        if input == "" {
            continue
        }

        cli.handleCommand(input)
    }

    if err := scanner.Err(); err != nil {
        fmt.Printf("❌ Error reading input: %v\n", err)
    }
}

// handleCommand processes user commands
func (cli *CLI) handleCommand(input string) {
    parts := strings.Fields(input)
    if len(parts) == 0 {
        return
    }

    command := strings.ToLower(parts[0])
    args := parts[1:]

    switch command {
    case "help", "h":
        cli.showHelp()
    case "exit", "quit", "q":
        cli.running = false
        fmt.Println("👋 Goodbye!")
    case "createaccount", "ca":
        cli.createAccount(args)
    case "importaccount", "ia":
        cli.importAccount(args)
    case "listaccounts", "la":
        cli.listAccounts()
    case "balance", "bal":
        cli.getBalance(args)
    case "stake", "s":
        cli.getStake(args)
    case "accountinfo", "ai":
        cli.getAccountInfo(args)
    case "send", "snd":
        cli.sendTransaction(args)
    case "staketokens", "st":
        cli.stakeTokens(args)
    case "setdefault", "sd":
        cli.setDefaultAccount(args)
    case "blockchain", "bc":
        cli.showBlockchainInfo()
    case "network", "net":
        cli.showNetworkInfo()
    default:
        fmt.Printf("❌ Unknown command: %s\n", command)
        fmt.Println("   Type 'help' for available commands")
    }
}

// showHelp displays available commands
func (cli *CLI) showHelp() {
    fmt.Println("")
    fmt.Println("📖 AVAILABLE COMMANDS:")
    fmt.Println("")
    fmt.Println("💰 Account Management:")
    fmt.Println("   createaccount (ca) [password]    - Create new account")
    fmt.Println("   importaccount (ia) [privkey] [pwd] - Import account")
    fmt.Println("   listaccounts (la)                - List all accounts")
    fmt.Println("   setdefault (sd) [address] [pwd]  - Set default account")
    fmt.Println("")
    fmt.Println("💸 Transactions:")
    fmt.Println("   send (snd) [to] [amount] [pwd]   - Send SELSI tokens")
    fmt.Println("   staketokens (st) [amount] [pwd]  - Stake SELSI tokens")
    fmt.Println("")
    fmt.Println("📊 Account Info:")
    fmt.Println("   balance (bal) [address]          - Check balance")
    fmt.Println("   stake (s) [address]              - Check stake")
    fmt.Println("   accountinfo (ai) [address]       - Full account info")
    fmt.Println("")
    fmt.Println("⛓️  Blockchain:")
    fmt.Println("   blockchain (bc)                  - Show blockchain info")
    fmt.Println("   network (net)                    - Show network info")
    fmt.Println("")
    fmt.Println("⚙️  System:")
    fmt.Println("   help (h)                         - Show this help")
    fmt.Println("   exit (q)                         - Exit CLI")
    fmt.Println("")
}

// createAccount creates a new account
func (cli *CLI) createAccount(args []string) {
    if len(args) < 1 {
        fmt.Println("❌ Usage: createaccount [password]")
        return
    }

    password := args[0]
    _, err := cli.wallet.CreateAccount(password)
    if err != nil {
        fmt.Printf("❌ Failed to create account: %v\n", err)
    }
}

// importAccount imports an account from private key
func (cli *CLI) importAccount(args []string) {
    if len(args) < 2 {
        fmt.Println("❌ Usage: importaccount [privatekey] [password]")
        return
    }

    privateKey := args[0]
    password := args[1]
    _, err := cli.wallet.ImportAccount(privateKey, password)
    if err != nil {
        fmt.Printf("❌ Failed to import account: %v\n", err)
    }
}

// listAccounts lists all accounts
func (cli *CLI) listAccounts() {
    addresses, err := cli.wallet.ListAccounts()
    if err != nil {
        fmt.Printf("❌ Failed to list accounts: %v\n", err)
        return
    }

    if len(addresses) == 0 {
        fmt.Println("📭 No accounts found")
        return
    }

    fmt.Println("")
    fmt.Println("👤 ACCOUNTS:")
    fmt.Println("------------")
    
    defaultAccount := cli.wallet.GetDefaultAccount()
    
    for i, address := range addresses {
        balance, _ := cli.wallet.GetBalance(address)
        stake, _ := cli.wallet.GetStake(address)
        
        status := ""
        if defaultAccount != nil && address == defaultAccount.Address {
            status = " [DEFAULT]"
        }
        
        fmt.Printf("%d. %s%s\n", i+1, hex.EncodeToString(address[:4]), status)
        fmt.Printf("   Balance: %s SELSI\n", new(big.Int).Div(balance, big.NewInt(1e18)))
        fmt.Printf("   Stake: %s SELSI\n", new(big.Int).Div(stake, big.NewInt(1e18)))
        fmt.Println()
    }
}

// getBalance shows account balance
func (cli *CLI) getBalance(args []string) {
    address, err := cli.getAddressFromArgs(args)
    if err != nil {
        fmt.Printf("❌ %v\n", err)
        return
    }

    balance, err := cli.wallet.GetBalance(address)
    if err != nil {
        fmt.Printf("❌ Failed to get balance: %v\n", err)
        return
    }

    fmt.Printf("💰 Balance: %s SELSI\n", new(big.Int).Div(balance, big.NewInt(1e18)))
}

// getStake shows account stake
func (cli *CLI) getStake(args []string) {
    address, err := cli.getAddressFromArgs(args)
    if err != nil {
        fmt.Printf("❌ %v\n", err)
        return
    }

    stake, err := cli.wallet.GetStake(address)
    if err != nil {
        fmt.Printf("❌ Failed to get stake: %v\n", err)
        return
    }

    fmt.Printf("🎯 Stake: %s SELSI\n", new(big.Int).Div(stake, big.NewInt(1e18)))
}

// getAccountInfo shows detailed account information
func (cli *CLI) getAccountInfo(args []string) {
    address, err := cli.getAddressFromArgs(args)
    if err != nil {
        fmt.Printf("❌ %v\n", err)
        return
    }

    info, err := cli.wallet.GetAccountInfo(address)
    if err != nil {
        fmt.Printf("❌ Failed to get account info: %v\n", err)
        return
    }

    fmt.Println("")
    fmt.Println("👤 ACCOUNT INFORMATION:")
    fmt.Println("----------------------")
    fmt.Printf("Address: %s\n", hex.EncodeToString(info.Address[:]))
    fmt.Printf("Balance: %s SELSI\n", new(big.Int).Div(info.Balance, big.NewInt(1e18)))
    fmt.Printf("Stake: %s SELSI\n", new(big.Int).Div(info.Stake, big.NewInt(1e18)))
    fmt.Printf("Nonce: %d\n", info.Nonce)
    fmt.Println("")
}

// sendTransaction sends SELSI tokens
func (cli *CLI) sendTransaction(args []string) {
    if len(args) < 3 {
        fmt.Println("❌ Usage: send [to_address] [amount] [password]")
        return
    }

    toHex := args[0]
    amountStr := args[1]
    password := args[2]

    // Parse amount
    amount, ok := new(big.Int).SetString(amountStr, 10)
    if !ok {
        fmt.Println("❌ Invalid amount")
        return
    }
    amount = new(big.Int).Mul(amount, big.NewInt(1e18)) // Convert to wei

    // Parse to address
    toBytes, err := hex.DecodeString(toHex)
    if err != nil || len(toBytes) != 20 {
        fmt.Println("❌ Invalid to address")
        return
    }
    var to types.Address
    copy(to[:], toBytes)

    // Get from address (default account)
    defaultAccount := cli.wallet.GetDefaultAccount()
    if defaultAccount == nil {
        fmt.Println("❌ No default account set")
        return
    }

    from := defaultAccount.Address

    tx, err := cli.wallet.SendTransaction(from, to, amount, password)
    if err != nil {
        fmt.Printf("❌ Failed to send transaction: %v\n", err)
        return
    }

    fmt.Printf("✅ Transaction created successfully!\n")
    fmt.Printf("   TX Hash: %s\n", "pending") // In real implementation, this would be the actual hash
}

// stakeTokens stakes SELSI tokens
func (cli *CLI) stakeTokens(args []string) {
    if len(args) < 2 {
        fmt.Println("❌ Usage: staketokens [amount] [password]")
        return
    }

    amountStr := args[0]
    password := args[1]

    // Parse amount
    amount, ok := new(big.Int).SetString(amountStr, 10)
    if !ok {
        fmt.Println("❌ Invalid amount")
        return
    }
    amount = new(big.Int).Mul(amount, big.NewInt(1e18)) // Convert to wei

    // Get from address (default account)
    defaultAccount := cli.wallet.GetDefaultAccount()
    if defaultAccount == nil {
        fmt.Println("❌ No default account set")
        return
    }

    from := defaultAccount.Address

    tx, err := cli.wallet.StakeTokens(from, amount, password)
    if err != nil {
        fmt.Printf("❌ Failed to stake tokens: %v\n", err)
        return
    }

    fmt.Printf("✅ Staking transaction created successfully!\n")
    fmt.Printf("   TX Hash: %s\n", "pending")
}

// setDefaultAccount sets the default account
func (cli *CLI) setDefaultAccount(args []string) {
    if len(args) < 2 {
        fmt.Println("❌ Usage: setdefault [address] [password]")
        return
    }

    addressHex := args[0]
    password := args[1]

    addressBytes, err := hex.DecodeString(addressHex)
    if err != nil || len(addressBytes) != 20 {
        fmt.Println("❌ Invalid address")
        return
    }
    var address types.Address
    copy(address[:], addressBytes)

    err = cli.wallet.SetDefaultAccount(address, password)
    if err != nil {
        fmt.Printf("❌ Failed to set default account: %v\n", err)
    }
}

// showBlockchainInfo shows blockchain information
func (cli *CLI) showBlockchainInfo() {
    currentBlock := cli.chain.GetCurrentBlock()
    
    fmt.Println("")
    fmt.Println("⛓️  BLOCKCHAIN INFO:")
    fmt.Println("------------------")
    fmt.Printf("Current Block: #%s\n", currentBlock.Header.Number)
    fmt.Printf("Total Blocks: %d\n", cli.chain.GetBlockCount())
    fmt.Printf("Chain ID: 769\n")
    fmt.Println("")
}

// showNetworkInfo shows network information
func (cli *CLI) showNetworkInfo() {
    fmt.Println("")
    fmt.Println("🌐 NETWORK INFO:")
    fmt.Println("---------------")
    fmt.Printf("Network: SelsiChain Mainnet\n")
    fmt.Printf("Consensus: Hybrid PoW/PoS\n")
    fmt.Printf("Status: Online\n")
    fmt.Println("")
}

// getAddressFromArgs extracts address from arguments
func (cli *CLI) getAddressFromArgs(args []string) (types.Address, error) {
    if len(args) > 0 {
        // Use provided address
        addressHex := args[0]
        addressBytes, err := hex.DecodeString(addressHex)
        if err != nil || len(addressBytes) != 20 {
            return types.Address{}, fmt.Errorf("invalid address")
        }
        var address types.Address
        copy(address[:], addressBytes)
        return address, nil
    }

    // Use default account
    defaultAccount := cli.wallet.GetDefaultAccount()
    if defaultAccount == nil {
        return types.Address{}, fmt.Errorf("no default account set and no address provided")
    }

    return defaultAccount.Address, nil
}
