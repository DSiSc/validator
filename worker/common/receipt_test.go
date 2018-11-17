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
	var except = types.Hash{
		0xe, 0x6e, 0xb5, 0x17, 0xf8, 0xf7, 0x5a, 0x3a, 0xf2, 0x28, 0xaf, 0xba, 0xce, 0xdc, 0x70, 0x24,
		0xca, 0x3c, 0x9e, 0xd, 0x6c, 0x2b, 0x3e, 0x51, 0xab, 0xf8, 0xc0, 0x4d, 0xcf, 0x23, 0xdf, 0xbf,
	}
	assert.Equal(except, receiptHash)
}
