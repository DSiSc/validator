package common

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/DSiSc/craft/types"
)

// TODO: Hash algorithm will support configurable later
// Sum returns the first 20 bytes of SHA256 of the bz.
func Sum(bz []byte) []byte {
	hash := sha256.Sum256(bz)
	return hash[:types.HashLength]
}

func TxHash(tx *types.Transaction) (hash types.Hash) {
	jsonByte, _ := json.Marshal(tx)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

func HeaderHash(block *types.Block) (hash types.Hash) {
	if *(new(types.Hash)) != block.HeaderHash {
		hash = block.HeaderHash
		return
	}
	header := block.Header
	jsonByte, _ := json.Marshal(header)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	block.HeaderHash = hash
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

func BlockHash(block *types.Block) (hash types.Hash) {
	jsonByte, _ := json.Marshal(block)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}

func ByteToHash(data []byte) (hash types.Hash) {
	copy(hash[:], data[:])
	return
}
