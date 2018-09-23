package common

import (
	"encoding/json"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/common"
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint64(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint64(1)
)

func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(root []byte, failed bool, cumulativeGasUsed uint64) *types.Receipt {
	r := &types.Receipt{PostState: CopyBytes(root), CumulativeGasUsed: cumulativeGasUsed}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}

// compute hash of receipt
func ReceiptHash(receipt *types.Receipt) (hash types.Hash) {
	jsonByte, _ := json.Marshal(receipt)
	sumByte := common.Sum(jsonByte)
	copy(hash[:], sumByte)
	return
}
