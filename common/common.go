package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
)

// TODO: Hash algorithm will support configurable later
// Sum returns the first 20 bytes of SHA256 of the bz.
func Sum(bz []byte) []byte {
	hash := sha256.Sum256(bz)
	return hash[:types.HashLength]
}

func TxHash(tx *types.Transaction) types.Hash {
	if hash := tx.Hash.Load(); hash != nil {
		return hash.(types.Hash)
	}
	hashData := types.TxData{
		AccountNonce: tx.Data.AccountNonce,
		Price: tx.Data.Price,
		GasLimit:tx.Data.GasLimit,
		Recipient:tx.Data.Recipient,
		Amount:tx.Data.Amount,
		Payload:tx.Data.Payload,
		V: tx.Data.V,
		R: tx.Data.R,
		S: tx.Data.S,
	}
	jsonByte, _ := json.Marshal(hashData)
	sumByte := Sum(jsonByte)
	var temp types.Hash
	copy(temp[:], sumByte)
	tx.Hash.Store(temp)
	return temp
}

func HeaderHash(block *types.Block) (hash types.Hash) {
	var defaultHash types.Hash
	if !bytes.Equal(block.HeaderHash[:], defaultHash[:]) {
		log.Info("block hash %v has exits.", block.HeaderHash)
		copy(hash[:], block.HeaderHash[:])
		return
	}
	jsonByte, _ := json.Marshal(block.Header)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

type RefAddress struct {
	Addr types.Address
}

func NewRefAddress(addr types.Address) *RefAddress {
	return &RefAddress{Addr: addr}
}

func (self *RefAddress) Address() types.Address {
	return self.Addr
}

func HashToByte(hash types.Hash) []byte {
	var value = make([]byte, len(hash))
	copy(value, hash[:])
	return value
}

func ByteToHash(data []byte) (hash types.Hash) {
	copy(hash[:], data[:])
	return
}

func HeaderDigest(header *types.Header) (hash types.Hash) {
	var defaultHash types.Hash
	if !bytes.Equal(header.MixDigest[:], defaultHash[:]) {
		log.Info("header hash %v has exits.", header.MixDigest)
		copy(hash[:], header.MixDigest[:])
		return
	}
	newHeader := digestHeader(header)
	jsonByte, _ := json.Marshal(newHeader)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

func digestHeader(header *types.Header) *types.Header {
	return &types.Header{
		ChainID:       header.ChainID,
		PrevBlockHash: header.PrevBlockHash,
		StateRoot:     header.StateRoot,
		TxRoot:        header.TxRoot,
		ReceiptsRoot:  header.ReceiptsRoot,
		Height:        header.Height,
		Timestamp:     header.Timestamp,
		Coinbase:      header.Coinbase,
	}
}
