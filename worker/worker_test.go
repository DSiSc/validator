package worker

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/tools"
	wallett "github.com/DSiSc/wallet/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestNewWorker(t *testing.T) {
	assert := assert.New(t)
	var worker = NewWorker(nil, nil)
	assert.NotNil(worker)
	assert.Nil(worker.block)
	assert.Nil(worker.chain)
}

func TestGetTxsRoot(t *testing.T) {
	var trxs = make([]*types.Transaction, 0)
	trx := new(types.Transaction)
	trxs = append(trxs, trx)
	hash := GetTxsRoot(trxs)
	assert.NotNil(t, hash)
}

var addressA = types.Address{
	0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
}

var addressB = tools.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")

func TestWorker_VerifyTrsSignature(t *testing.T) {
	key, _ := wallett.DefaultTestKey()
	mockTrx := &types.Transaction{
		Data: types.TxData{
			AccountNonce: uint64(0),
			Price:        new(big.Int),
			Recipient:    &addressA,
			From:         &addressB,
			Amount:       new(big.Int),
			Payload:      addressB[:10],
		},
	}
	mockTransaction, _ := wallett.SignTx(mockTrx, new(wallett.FrontierSigner), key)
	worker := NewWorker(nil, nil)
	ok := worker.VerifyTrsSignature(mockTransaction)
	assert.Equal(t, true, ok)
}
