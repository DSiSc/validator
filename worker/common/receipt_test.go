package common

import (
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
