package common

import (
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewReceipt(t *testing.T) {
	assert := assert.New(t)
	recept := NewReceipt(nil, true, uint64(0))
	assert.NotNil(recept)
	assert.Equal(recept.Status, ReceiptStatusFailed)
	recept = NewReceipt(nil, false, uint64(0))
	assert.Equal(recept.Status, ReceiptStatusSuccessful)
}

func TestReceipt_ReceiptHash(t *testing.T) {
	assert := assert.New(t)
	receipt := NewReceipt(nil, false, uint64(0))
	assert.Equal(receipt.Status, ReceiptStatusSuccessful)
	receiptHash := ReceiptHash(receipt)
	var except = types.Hash{0x1, 0x8, 0x22, 0xb1, 0x12, 0x1, 0xd3, 0x70, 0x7d, 0x48, 0xc5, 0x4b, 0x2, 0xf4, 0x82, 0x7d, 0x95, 0x4e, 0xfd, 0x25, 0x6a, 0xc, 0xaa, 0xf5, 0x60, 0x62, 0x39, 0x6, 0x52, 0x3a, 0x2a, 0xa5}

	assert.Equal(except, receiptHash)
}
