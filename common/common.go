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

func TxHash(tx *types.Transaction) (hash types.Hash) {
	jsonByte, _ := json.Marshal(tx)
	sumByte := Sum(jsonByte)
	copy(hash[:], sumByte)
	return
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
