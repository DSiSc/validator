package worker

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/repository"
	"github.com/DSiSc/validator/worker/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
	"time"
)

var MockHash = types.Hash{
	0x1d, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

var MockBlock = &types.Block{
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

var from = &types.Address{
	0xb1, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
}

var to = &types.Address{
	0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
}

var contractAddress = types.Address{
	0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x16,
}

func mockTrx() *types.Transaction {
	return &types.Transaction{
		Data: types.TxData{
			AccountNonce: 0,
			Price:        new(big.Int).SetUint64(10),
			GasLimit:     100,
			Recipient:    to,
			From:         from,
			Amount:       new(big.Int).SetUint64(50),
			Payload:      to[:10],
		},
	}
}

var state *StateTransition

func TestNewStateTransition(t *testing.T) {
	evmNg := &evm.EVM{
		StateDB: nil,
	}
	var gp = common.GasPool(6)
	state = NewStateTransition(evmNg, mockTrx(), &gp)
	assert.NotNil(t, state.evm)
	assert.NotNil(t, state.tx)
}

func TestStateTransition_TransitionDb(t *testing.T) {
	// test carets contract
	evmNg := &evm.EVM{
		StateDB: nil,
	}
	var gp = common.GasPool(10000)
	state = NewStateTransition(evmNg, mockTrx(), &gp)
	state.tx.Data.Recipient = nil
	var evmd *evm.EVM
	monkey.PatchInstanceMethod(reflect.TypeOf(evmd), "Create", func(*evm.EVM, evm.ContractRef, []byte, uint64, *big.Int) ([]byte, types.Address, uint64, error) {
		return to[:10], contractAddress, 0, evm.ErrInsufficientBalance
	})
	var bc *repository.Repository
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetBalance", func(*repository.Repository, types.Address) *big.Int {
		return new(big.Int).SetUint64(1000)
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SubBalance", func(*repository.Repository, types.Address, *big.Int) {
		return
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetNonce", func(*repository.Repository, types.Address) uint64 {
		return 0
	})
	ret, used, ok, err, address := state.TransitionDb()
	assert.Equal(t, err, evm.ErrInsufficientBalance)
	assert.Equal(t, ok, false)
	assert.Equal(t, uint64(0), used)
	assert.Equal(t, to[:10], ret)
	assert.Equal(t, contractAddress, address)
	monkey.UnpatchAll()
}

func TestStateTransition_TransitionDb1(t *testing.T) {
	defer monkey.UnpatchAll()
	// test transfer token
	evmNg := &evm.EVM{
		StateDB: nil,
	}
	var gp = common.GasPool(10000)
	state := NewStateTransition(evmNg, mockTrx(), &gp)
	var evmd *evm.EVM
	monkey.PatchInstanceMethod(reflect.TypeOf(evmd), "Create", func(*evm.EVM, evm.ContractRef, []byte, uint64, *big.Int) ([]byte, types.Address, uint64, error) {
		return to[:10], contractAddress, 0, evm.ErrInsufficientBalance
	})
	var bc *repository.Repository
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetBalance", func(*repository.Repository, types.Address) *big.Int {
		return new(big.Int).SetUint64(1000)
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SubBalance", func(*repository.Repository, types.Address, *big.Int) {
		return
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetNonce", func(*repository.Repository, types.Address) uint64 {
		return 0
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SetNonce", func(b *repository.Repository, a types.Address, n uint64) {
		assert.Equal(t, uint64(1), n)
		return
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(evmd), "Call", func(*evm.EVM, evm.ContractRef, types.Address, []byte, uint64, *big.Int) ([]byte, uint64, error) {
		return []byte{0}, 0, nil
	})
	ret, used, ok, err, _ := state.TransitionDb()
	assert.Equal(t, nil, err)
	assert.Equal(t, ok, false)
	assert.Equal(t, uint64(0), used)
	assert.Equal(t, []byte{0}, ret)

	state.nonce = 100
	ret, used, ok, err, _ = state.TransitionDb()
	assert.NotNil(t, err)
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetNonce", func(*repository.Repository, types.Address) uint64 {
		return 3
	})
	ret, used, ok, err, _ = state.TransitionDb()
	assert.NotNil(t, err)
}
