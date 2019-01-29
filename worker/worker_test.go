package worker

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools"
	workerc "github.com/DSiSc/validator/worker/common"
	walletc "github.com/DSiSc/wallet/common"
	wallett "github.com/DSiSc/wallet/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
)

func TestNewWorker(t *testing.T) {
	assert := assert.New(t)
	var worker = NewWorker(nil, nil, false)
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

var addressC = walletc.Address{
	0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x16,
}

var addressNew = types.Address{
	0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x17,
}

var addressB = tools.HexToAddress("0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b")

var mockHash = types.Hash{
	0x1d, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

var mockHash1 = types.Hash{
	0x1e, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

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
	worker := NewWorker(nil, nil, false)
	ok := worker.VerifyTrsSignature(mockTransaction)
	assert.Equal(t, true, ok)

	monkey.Patch(wallett.Sender, func(wallett.Signer, *types.Transaction) (walletc.Address, error) {
		return addressC, fmt.Errorf("unknown signer")
	})
	ok = worker.VerifyTrsSignature(mockTransaction)
	assert.Equal(t, false, ok)

	monkey.Patch(wallett.Sender, func(wallett.Signer, *types.Transaction) (walletc.Address, error) {
		return addressC, nil
	})
	ok = worker.VerifyTrsSignature(mockTransaction)
	assert.Equal(t, false, ok)
	monkey.UnpatchAll()
}

func TestWorker_VerifyBlock(t *testing.T) {
	assert := assert.New(t)
	var blockChain *blockchain.BlockChain
	var mockBlock = &types.Block{
		Header: &types.Header{
			ChainID: uint64(1),
			Height:  uint64(1),
		},
	}
	worker := NewWorker(nil, mockBlock, false)

	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		return &types.Block{
			Header: &types.Header{
				ChainID: uint64(0),
			},
		}
	})
	err := worker.VerifyBlock()
	assert.NotNil(err, "chain id not consistent")

	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		return &types.Block{
			Header: &types.Header{
				ChainID: uint64(1),
			},
			HeaderHash: mockHash,
		}
	})
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block pre block hash not consistent")

	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		return &types.Block{
			Header: &types.Header{
				ChainID: uint64(1),
			},
		}
	})
	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlockHeight", func(*blockchain.BlockChain) uint64 {
		return 1
	})
	//mockBlock.Header.ChainID = uint64(0)
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block height not consistent")

	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlockHeight", func(*blockchain.BlockChain) uint64 {
		return 0
	})
	worker.block.Header.TxRoot = MockHash
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block txroot hash not consistent")

	monkey.Patch(common.HeaderHash, func(*types.Block) types.Hash {
		return MockHash
	})
	var tmp types.Hash
	worker.block.Header.TxRoot = tmp
	worker.block.HeaderHash = mockHash1
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block header hash not consistent")

	monkey.Patch(common.HeaderHash, func(*types.Block) types.Hash {
		var tmp types.Hash
		return tmp
	})
	worker.block.Header.ReceiptsRoot = MockHash
	worker.block.HeaderHash = common.HeaderHash(worker.block)
	err = worker.VerifyBlock()
	assert.NotNil(err, "Receipts hash not consistent")

	worker.block.Header.ReceiptsRoot = tmp
	err = worker.VerifyBlock()
	assert.Nil(err)
	monkey.UnpatchAll()

}

func TestWorker_VerifyTransaction(t *testing.T) {
	assert := assert.New(t)
	worker := NewWorker(nil, nil, false)

	monkey.Patch(evm.NewEVMContext, func(types.Transaction, *types.Header, *blockchain.BlockChain, types.Address) evm.Context {
		return evm.Context{
			GasLimit: uint64(65536),
		}
	})
	monkey.Patch(ApplyTransaction, func(*evm.EVM, *types.Transaction, *workerc.GasPool) ([]byte, uint64, bool, error) {
		return addressA[:10], uint64(0), false, fmt.Errorf("Apply failed.")
	})
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
	receipit, gas, err := worker.VerifyTransaction(addressA, nil, nil, mockTrx, nil)
	assert.Equal(err, fmt.Errorf("Apply failed."))
	assert.Nil(receipit)
	assert.Equal(uint64(0), gas)

	monkey.Patch(ApplyTransaction, func(*evm.EVM, *types.Transaction, *workerc.GasPool) ([]byte, uint64, bool, error) {
		return addressA[:10], uint64(10), true, nil
	})
	var blockChain *blockchain.BlockChain
	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "IntermediateRoot", func(*blockchain.BlockChain, bool) types.Hash {
		return MockHash
	})
	monkey.Patch(crypto.CreateAddress, func(types.Address, uint64) types.Address {
		return addressNew
	})
	monkey.Patch(evm.NewEVMContext, func(types.Transaction, *types.Header, *blockchain.BlockChain, types.Address) evm.Context {
		return evm.Context{
			GasLimit: uint64(65536),
		}
	})
	mockTrx.Data.Recipient = nil
	var usedGas = uint64(10)
	receipit, gas, err = worker.VerifyTransaction(addressA, nil, nil, mockTrx, &usedGas)
	assert.Equal(addressNew, receipit.ContractAddress)
	assert.Equal(uint64(10), gas)
	monkey.UnpatchAll()
}

func TestWorker_GetReceipts(t *testing.T) {
	assert := assert.New(t)
	worker := NewWorker(nil, nil, false)
	receipts := worker.GetReceipts()
	assert.Equal(len(receipts), len(worker.receipts))
}
