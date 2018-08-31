package worker

import (
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/evm-NG"
	"github.com/DSiSc/evm-NG/common/crypto"
	"github.com/DSiSc/validator/tools"
)

type Worker struct {
	block    *types.Block
	chain    *blockchain.BlockChain
	receipts types.Receipts
	logs     []*types.Log
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
		txHash = append(txHash, t.Hash())
	}
	txRoot := tools.ComputeMerkleRoot(txHash)
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
	if self.block.Header.PrevBlockHash != currentBlock.Hash() {
		return fmt.Errorf("Wrong Block.Header.PrevBlockHash. Expected %v, got %v",
			currentBlock.Hash(), self.block.Header.PrevBlockHash)
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
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = self.block.Header
		allLogs  []*types.Log
		gp       = new(types.GasPool).AddGas(uint64(65536))
	)

	var from types.Address
	var sign types.Signer
	// 5. verify every transactions in the block by evm
	for i, tx := range self.block.Transactions {
		self.chain.Prepare(tx.Hash(), self.block.Hash(), i)
		from, _ = types.Sender(sign, tx)
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
	gp *types.GasPool,
	header *types.Header,
	tx *types.Transaction,
	usedGas *uint64) (*types.Receipt, uint64, error) {

	var signer types.Signer
	msg, err := tx.AsMessage(signer)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed to make message.")
	}

	context := evm.NewEVMContext(msg, header, self.chain, author)
	evmEnv := evm.NewEVM(context, self.chain)
	_, gas, failed, err := ApplyMessage(evmEnv, msg, gp)
	if err != nil {
		return nil, 0, err
	}

	root := self.chain.IntermediateRoot(false).Bytes()
	*usedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing wether the root touch-delete accounts.
	receipt := types.NewReceipt(root, failed, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas
	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(evmEnv.Context.Origin, uint64(0))
	}
	// Set the receipt logs and create a bloom for filtering
	// receipt.Logs = self.chain.GetLogs(tx.Hash())
	// receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}
