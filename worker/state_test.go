package worker

import (
	"encoding/hex"
	"encoding/json"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/evm-NG/common/math"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/repository"
	"github.com/DSiSc/validator/worker/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
	"time"
	"fmt"
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

var author = types.Address{
	0xb1, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xc2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
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

var wasmContractAddress = types.Address{
	0x79, 0x46, 0x28, 0x8, 0xf6, 0xd1, 0xa6, 0x42, 0x81, 0xd, 0x96, 0xa1, 0xfb, 0x67, 0x5c, 0x33, 0xcf, 0x60, 0xc8, 0x65,
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
	bc := &repository.Repository{}
	var gp = common.GasPool(6)
	state = NewStateTransition(author, MockBlock.Header, bc, mockTrx(), &gp)
	assert.NotNil(t, state)
	assert.NotNil(t, state.tx)
}

func TestStateTransition_TransitionDb(t *testing.T) {
	// test carets contract
	bc := &repository.Repository{}
	var gp = common.GasPool(10000)
	state = NewStateTransition(author, MockBlock.Header, bc, mockTrx(), &gp)
	state.tx.Data.Recipient = nil
	var evmd *evm.EVM
	monkey.PatchInstanceMethod(reflect.TypeOf(evmd), "Create", func(*evm.EVM, evm.ContractRef, []byte, uint64, *big.Int) ([]byte, types.Address, uint64, error) {
		return to[:10], contractAddress, 0, evm.ErrInsufficientBalance
	})
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
	bc := &repository.Repository{}
	var gp = common.GasPool(10000)
	state = NewStateTransition(author, MockBlock.Header, bc, mockTrx(), &gp)
	var evmd *evm.EVM
	monkey.PatchInstanceMethod(reflect.TypeOf(evmd), "Create", func(*evm.EVM, evm.ContractRef, []byte, uint64, *big.Int) ([]byte, types.Address, uint64, error) {
		return to[:10], contractAddress, 0, evm.ErrInsufficientBalance
	})
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
		return []byte{0}, math.MaxUint64, nil
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetCode", func(*repository.Repository, types.Address) []byte {
		return []byte{}
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

// test create wasm contract
func TestStateTransition_TransitionDb2(t *testing.T) {
	// test carets contract
	bc := &repository.Repository{}
	var gp = common.GasPool(10000)
	tx := mockTrx()
	code, _ := hex.DecodeString("0061736d0100000001070160027f7f017f03020100070801046961646400000a09010700200020016a0b")
	tx.Data.Payload = code
	state = NewStateTransition(author, MockBlock.Header, bc, tx, &gp)
	state.tx.Data.Recipient = nil
	var evmd *evm.EVM
	monkey.PatchInstanceMethod(reflect.TypeOf(evmd), "Create", func(*evm.EVM, evm.ContractRef, []byte, uint64, *big.Int) ([]byte, types.Address, uint64, error) {
		return to[:10], contractAddress, 0, evm.ErrInsufficientBalance
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetBalance", func(*repository.Repository, types.Address) *big.Int {
		return new(big.Int).SetUint64(1000)
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SubBalance", func(*repository.Repository, types.Address, *big.Int) {
		return
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetNonce", func(*repository.Repository, types.Address) uint64 {
		return 0
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "AddBalance", func(*repository.Repository, types.Address, *big.Int) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SetNonce", func(*repository.Repository, types.Address, uint64) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "CreateAccount", func(*repository.Repository, types.Address) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SetCode", func(*repository.Repository, types.Address, []byte) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetCodeHash", func(*repository.Repository, types.Address) types.Hash {
		return types.Hash{}
	})
	_, used, ok, err, address := state.TransitionDb()
	assert.Nil(t, err)
	assert.Equal(t, ok, false)
	assert.Equal(t, uint64(0), used)
	assert.Equal(t, wasmContractAddress, address)
	monkey.UnpatchAll()
}

// test call wasm contract
func TestStateTransition_TransitionDb3(t *testing.T) {
	defer monkey.UnpatchAll()
	// test carets contract
	bc := &repository.Repository{}
	var gp = common.GasPool(10000)
	tx := mockTrx()
	tx.Data.Recipient = &wasmContractAddress
	code, _ := hex.DecodeString("0061736d01000000018c808080000260017f017f60027f7f017f028e808080000103656e76066d616c6c6f6300000382808080000101048480808000017000000583808080000100010681808080000007938080800002066d656d6f7279020006696e766f6b6500010a998080800001938080800001017f41021000220241c8d2013b000020020b")
	tx.Data.Payload, _ = json.Marshal([]string{"Hi", "Bob"})
	fmt.Printf("%x", tx.Data.Payload)
	state = NewStateTransition(author, MockBlock.Header, bc, tx, &gp)
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetBalance", func(*repository.Repository, types.Address) *big.Int {
		return new(big.Int).SetUint64(1000)
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SubBalance", func(*repository.Repository, types.Address, *big.Int) {
		return
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetNonce", func(*repository.Repository, types.Address) uint64 {
		return 0
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "AddBalance", func(*repository.Repository, types.Address, *big.Int) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SetNonce", func(*repository.Repository, types.Address, uint64) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "CreateAccount", func(*repository.Repository, types.Address) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetCode", func(*repository.Repository, types.Address) []byte {
		return code
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "SetCode", func(*repository.Repository, types.Address, []byte) {
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(bc), "GetCodeHash", func(*repository.Repository, types.Address) types.Hash {
		return types.Hash{}
	})
	ret, used, ok, err, address := state.TransitionDb()
	assert.Nil(t, err)
	assert.Equal(t, []byte{'H', 'i'}, ret)
	assert.Equal(t, ok, false)
	assert.Equal(t, uint64(0), used)
	assert.Equal(t, types.Address{}, address)
	monkey.UnpatchAll()
}
