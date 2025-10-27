package hash

import (
    "crypto/sha256"
    "github.com/selsichain/selsichain-core/core/types"
)

func CalculateBlockHash(header *types.Header) types.Hash {
    data := serializeHeader(header)
    hash := sha256.Sum256(data)
    return types.Hash(hash)
}

func CalculateTransactionHash(tx *types.Transaction) types.Hash {
    data := serializeTransaction(tx)
    hash := sha256.Sum256(data)
    return types.Hash(hash)
}

func serializeHeader(header *types.Header) []byte {
    var data []byte
    data = append(data, header.ParentHash[:]...)
    data = append(data, header.Coinbase[:]...)
    if header.Number != nil {
        data = append(data, header.Number.Bytes()...)
    }
    data = append(data, byte(header.Time))
    if header.Difficulty != nil {
        data = append(data, header.Difficulty.Bytes()...)
    }
    data = append(data, header.Nonce[:]...)
    data = append(data, header.Validator[:]...)
    return data
}

func serializeTransaction(tx *types.Transaction) []byte {
    var data []byte
    data = append(data, byte(tx.Nonce))
    
    // Handle To address (bisa nil untuk contract creation)
    if tx.To != nil {
        data = append(data, tx.To[:]...)
    } else {
        // Empty address untuk contract creation
        data = append(data, make([]byte, 20)...)
    }
    
    if tx.Value != nil {
        data = append(data, tx.Value.Bytes()...)
    }
    
    // Use Data field (dulunya Input)
    data = append(data, tx.Data...)
    
    return data
}
