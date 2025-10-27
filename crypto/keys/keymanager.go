package keys

import (
    "crypto/ecdsa"
    "crypto/rand"
    "encoding/hex"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"

    "github.com/ethereum/go-ethereum/crypto"
    "github.com/selsichain/selsichain-core/core/types"
)

// KeyManager manages cryptographic keys
type KeyManager struct {
    keyDir string
}

// KeyPair represents a public-private key pair
type KeyPair struct {
    PrivateKey *ecdsa.PrivateKey
    PublicKey  *ecdsa.PublicKey
    Address    types.Address
}

// NewKeyManager creates a new key manager
func NewKeyManager(keyDir string) *KeyManager {
    // Create directory if it doesn't exist
    os.MkdirAll(keyDir, 0700)
    
    return &KeyManager{
        keyDir: keyDir,
    }
}

// GenerateKey generates a new key pair
func (km *KeyManager) GenerateKey() (*KeyPair, error) {
    privateKey, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
    if err != nil {
        return nil, fmt.Errorf("failed to generate private key: %w", err)
    }

    return km.privateKeyToKeyPair(privateKey), nil
}

// ImportPrivateKey imports a private key from hex string
func (km *KeyManager) ImportPrivateKey(privateKeyHex string) (*KeyPair, error) {
    privateKeyBytes, err := hex.DecodeString(privateKeyHex)
    if err != nil {
        return nil, fmt.Errorf("invalid private key hex: %w", err)
    }

    privateKey, err := crypto.ToECDSA(privateKeyBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %w", err)
    }

    return km.privateKeyToKeyPair(privateKey), nil
}

// SaveKey saves a key pair to disk
func (km *KeyManager) SaveKey(keyPair *KeyPair, password string) error {
    privateKeyBytes := crypto.FromECDSA(keyPair.PrivateKey)
    
    // Simple encryption (in production, use proper encryption)
    encryptedData := km.simpleEncrypt(privateKeyBytes, password)
    
    filename := filepath.Join(km.keyDir, hex.EncodeToString(keyPair.Address[:])+".key")
    return ioutil.WriteFile(filename, encryptedData, 0600)
}

// LoadKey loads a key pair from disk
func (km *KeyManager) LoadKey(address types.Address, password string) (*KeyPair, error) {
    filename := filepath.Join(km.keyDir, hex.EncodeToString(address[:])+".key")
    encryptedData, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to read key file: %w", err)
    }

    privateKeyBytes := km.simpleDecrypt(encryptedData, password)
    privateKey, err := crypto.ToECDSA(privateKeyBytes)
    if err != nil {
        return nil, fmt.Errorf("failed to parse private key: %w", err)
    }

    return km.privateKeyToKeyPair(privateKey), nil
}

// ListKeys lists all saved keys
func (km *KeyManager) ListKeys() ([]types.Address, error) {
    files, err := ioutil.ReadDir(km.keyDir)
    if err != nil {
        return nil, fmt.Errorf("failed to read key directory: %w", err)
    }

    var addresses []types.Address
    for _, file := range files {
        if filepath.Ext(file.Name()) == ".key" {
            hexStr := file.Name()[:len(file.Name())-4] // Remove .key extension
            addressBytes, err := hex.DecodeString(hexStr)
            if err != nil {
                continue
            }
            
            var address types.Address
            copy(address[:], addressBytes)
            addresses = append(addresses, address)
        }
    }

    return addresses, nil
}

// privateKeyToKeyPair converts private key to key pair
func (km *KeyManager) privateKeyToKeyPair(privateKey *ecdsa.PrivateKey) *KeyPair {
    publicKey := &privateKey.PublicKey
    address := km.publicKeyToAddress(publicKey)
    
    return &KeyPair{
        PrivateKey: privateKey,
        PublicKey:  publicKey,
        Address:    address,
    }
}

// publicKeyToAddress converts public key to address
func (km *KeyManager) publicKeyToAddress(publicKey *ecdsa.PublicKey) types.Address {
    publicKeyBytes := crypto.FromECDSAPub(publicKey)
    // Ethereum-style address: last 20 bytes of keccak256 hash
    hash := crypto.Keccak256(publicKeyBytes[1:]) // Remove prefix
    var address types.Address
    copy(address[:], hash[12:]) // Last 20 bytes
    return address
}

// Simple encryption for demo (in production, use proper encryption)
func (km *KeyManager) simpleEncrypt(data []byte, password string) []byte {
    // XOR with password bytes (very basic, for demo only)
    result := make([]byte, len(data))
    passwordBytes := []byte(password)
    
    for i := 0; i < len(data); i++ {
        result[i] = data[i] ^ passwordBytes[i%len(passwordBytes)]
    }
    
    return result
}

// Simple decryption for demo
func (km *KeyManager) simpleDecrypt(data []byte, password string) []byte {
    // XOR is symmetric, so encryption and decryption are the same
    return km.simpleEncrypt(data, password)
}

// GetPrivateKeyHex returns private key as hex string
func (kp *KeyPair) GetPrivateKeyHex() string {
    return hex.EncodeToString(crypto.FromECDSA(kp.PrivateKey))
}

// GetAddressHex returns address as hex string
func (kp *KeyPair) GetAddressHex() string {
    return hex.EncodeToString(kp.Address[:])
}
