package validator

import (
	"github.com/DSiSc/craft/types"
	account2 "github.com/DSiSc/validator/tools/account"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var validator *Validator

var MockHash = types.Hash{
	0x1d, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

func MockBlock() *types.Block {
	return &types.Block{
		Header: &types.Header{
			ChainID:       1,
			PrevBlockHash: MockHash,
			StateRoot:     MockHash,
			TxRoot:        MockHash,
			ReceiptsRoot:  MockHash,
			Height:        1,
			Timestamp:     uint64(time.Date(2018, time.August, 28, 0, 0, 0, 0, time.UTC).Unix()),
			MixDigest:     MockHash,
		},
		Transactions: make([]*types.Transaction, 0),
	}
}

func TestNewValidator(t *testing.T) {
	var account account2.Account
	validator = NewValidator(&account)
	assert.NotNil(t, validator)
}

func TestValidateBlock(t *testing.T) {
	header, err := validator.ValidateBlock(MockBlock())
	assert.Nil(t, header)
	assert.NotNil(t, err)
}
