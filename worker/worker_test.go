package worker

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/monkey"
	"github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools"
	walletc "github.com/DSiSc/wallet/common"
	wallett "github.com/DSiSc/wallet/core/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
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

var addressC = walletc.Address{
	0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x16,
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

	exceptErr := fmt.Errorf("Unknown signer")
	monkey.Patch(wallett.Sender, func(wallett.Signer, *types.Transaction) (walletc.Address, error) {
		return addressC, exceptErr
	})
	ok = worker.VerifyTrsSignature(mockTransaction)
	assert.Equal(t, false, ok)

	monkey.Patch(wallett.Sender, func(wallett.Signer, *types.Transaction) (walletc.Address, error) {
		return addressC, nil
	})
	ok = worker.VerifyTrsSignature(mockTransaction)
	assert.Equal(t, false, ok)
	monkey.Unpatch(wallett.Sender)
}

func TestWorker_VerifyBlock(t *testing.T) {
	assert := assert.New(t)
	var blockChain *blockchain.BlockChain
	var mockBlock = &types.Block{
		Header: &types.Header{
			ChainID: uint64(1),
		},
	}
	worker := NewWorker(nil, mockBlock)

	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlock", func(*blockchain.BlockChain) *types.Block {
		return mockBlock
	})
	err := worker.VerifyBlock()
	assert.NotNil(err, "chain id not consistent")

	mockBlock.Header.ChainID = uint64(0)
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block pre block hash not consistent")

	monkey.Patch(common.BlockHash, func(*types.Block) types.Hash {
		return MockHash
	})
	worker.block.Header.PrevBlockHash = MockHash
	worker.block.Header.Height = uint64(0)
	monkey.PatchInstanceMethod(reflect.TypeOf(blockChain), "GetCurrentBlockHeight", func(*blockchain.BlockChain) uint64 {
		return 0
	})
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block height not consistent")

	worker.block.Header.Height = uint64(1)
	worker.block.Header.TxRoot = MockHash
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block txroot hash not consistent")

	worker.block.Header.TxRoot = GetTxsRoot(worker.block.Transactions)
	worker.block.HeaderHash = MockHash
	err = worker.VerifyBlock()
	assert.NotNil(err, "Block header hash not consistent")

	monkey.Patch(common.HeaderHash, func(*types.Block) types.Hash {
		return MockHash
	})
	worker.block.Header.ReceiptsRoot = MockHash
	err = worker.VerifyBlock()
	assert.NotNil(err, "Receipts hash not consistent")

	var temp types.Hash
	worker.block.Header.ReceiptsRoot = temp
	worker.block.HeaderHash = common.HeaderHash(worker.block)
	err = worker.VerifyBlock()
	assert.Nil(err)
}
