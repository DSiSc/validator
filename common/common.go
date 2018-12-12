package common

import (
	"encoding/json"
	gconf "github.com/DSiSc/craft/config"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto/sha3"
)

// Sum returns the first 32 bytes of hash of the bz.
func Sum(bz []byte) []byte {
	var alg string
	if value, ok := gconf.GlobalConfig.Load(gconf.HashAlgName); ok {
		alg = value.(string)
	} else {
		alg = "SHA256"
	}
	hasher := sha3.NewHashByAlgName(alg)
	hasher.Write(bz)
	hash := hasher.Sum(nil)
	return hash[:types.HashLength]
}

func TxHash(tx *types.Transaction) types.Hash {
	if hash := tx.Hash.Load(); hash != nil {
		return hash.(types.Hash)
	}
	hashData := types.TxData{
		AccountNonce: tx.Data.AccountNonce,
		Price:        tx.Data.Price,
		GasLimit:     tx.Data.GasLimit,
		Recipient:    tx.Data.Recipient,
		Amount:       tx.Data.Amount,
		Payload:      tx.Data.Payload,
		V:            tx.Data.V,
		R:            tx.Data.R,
		S:            tx.Data.S,
	}
	jsonByte, _ := json.Marshal(hashData)
	sumByte := Sum(jsonByte)
	var temp types.Hash
	copy(temp[:], sumByte)
	tx.Hash.Store(temp)
	return temp
}

func HeaderHash(block *types.Block) types.Hash {
	var headerHash types.Hash
	if !(block.HeaderHash == types.Hash{}) {
		return block.HeaderHash
	}
	jsonByte, _ := json.Marshal(block.Header)
	sumByte := Sum(jsonByte)
	copy(headerHash[:], sumByte)
	return headerHash
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
	if !(header.MixDigest == types.Hash{}) {
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
		CoinBase:      header.CoinBase,
	}
}
