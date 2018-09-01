package worker

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/evm-NG/common/crypto"
	vcommon "github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools/merkle_tree"
	"github.com/DSiSc/validator/tools/signature"
	"github.com/DSiSc/validator/worker/common"
	//"github.com/basechain/core"
)

type Worker struct {
	block    *types.Block
	chain    *blockchain.BlockChain
	receipts common.Receipts
	logs     []*common.Log
}

func NewWorker(chain *blockchain.BlockChain, block *types.Block) *Worker {
	return &Worker{
		block: block,
		chain: chain,
	}
}

func GetTxsRoot(txs []*types.Transaction) types.Hash {
	txHash := make([]types.Hash, 0, len(txs))
	for _, t := range txs {
		txHash = append(txHash, vcommon.TxHash(t))
	}
	txRoot := merkle_tree.ComputeMerkleRoot(txHash)
	return txRoot
}

func (self *Worker) VerifyBlock() error {
	// TODO: Verify the block attributyes
	// 1. chainID
	currentBlock := self.chain.GetCurrentBlock()
	if self.block.Header.ChainID != currentBlock.Header.ChainID {
		return fmt.Errorf("Wrong Block.Header.ChainID. Expected %v, got %v",
			currentBlock.Header.ChainID, self.block.Header.ChainID)
	}

	// 2. hash
	if self.block.Header.PrevBlockHash != vcommon.BlockHash(currentBlock) {
		return fmt.Errorf("Wrong Block.Header.PrevBlockHash. Expected %v, got %v",
			vcommon.BlockHash(currentBlock), self.block.Header.PrevBlockHash)
	}

	// 3. height
	if self.block.Header.Height != self.chain.GetCurrentBlockHeight()+1 {
		return fmt.Errorf("Wrong Block.Header.Height. Expected %v, got %v",
			self.chain.GetCurrentBlockHeight()+1, self.block.Header.Height)
	}
	// 4. txhash
	txsHash := GetTxsRoot(self.block.Transactions)
	if self.block.Header.TxRoot != txsHash {
		return fmt.Errorf("Wrong Block.Header.TxRoot. Expected %v, got %v",
			txsHash, self.block.Header.TxRoot)
	}

	var (
		receipts common.Receipts
		usedGas  = new(uint64)
		header   = self.block.Header
		allLogs  []*common.Log
		gp       = new(common.GasPool).AddGas(uint64(65536))
	)

	var from types.Address
	var sign signature.TxSigner
	// 5. verify every transactions in the block by evm
	for i, tx := range self.block.Transactions {
		self.chain.Prepare(vcommon.TxHash(tx), vcommon.BlockHash(self.block), i)
		from, _ = signature.Sender(sign, tx)
		receipt, _, err := self.VerifyTransaction(from, gp, header, tx, usedGas)
		if err != nil {
			return err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	self.receipts = receipts
	self.logs = allLogs

	return nil
}

func (self *Worker) VerifyTransaction(
	author types.Address,
	gp *common.GasPool,
	header *types.Header,
	tx *types.Transaction,
	usedGas *uint64) (*common.Receipt, uint64, error) {

	context := evm.NewEVMContext(*tx, header, self.chain, author)
	evmEnv := evm.NewEVM(context, self.chain)
	_, gas, failed, err := ApplyTransaction(evmEnv, tx, gp)
	if err != nil {
		return nil, 0, err
	}

	root := self.chain.IntermediateRoot(false)
	*usedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	receipt := common.NewReceipt(vcommon.HashToByte(root), failed, *usedGas)
	receipt.TxHash = vcommon.TxHash(tx)
	receipt.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	if tx.Data.Recipient == nil {
		receipt.ContractAddress = crypto.CreateAddress(evmEnv.Context.Origin, uint64(0))
	}
	// Set the receipt logs and create a bloom for filtering
	// receipt.Logs = self.chain.GetLogs(tx.Hash())
	// receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}
