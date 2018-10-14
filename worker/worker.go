package worker

import (
	"bytes"
	"fmt"
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/crypto-suite/crypto"
	"github.com/DSiSc/evm-NG"
	vcommon "github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools/merkle_tree"
	"github.com/DSiSc/validator/worker/common"
	wallett "github.com/DSiSc/wallet/core/types"
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
		txHash = append(txHash, vcommon.TxHash(t))
	}
	txRoot := merkle_tree.ComputeMerkleRoot(txHash)
	return txRoot
}

func (self *Worker) VerifyBlock() error {
	// 1. chainID
	currentBlock := self.chain.GetCurrentBlock()
	if self.block.Header.ChainID != currentBlock.Header.ChainID {
		return fmt.Errorf("Wrong Block.Header.ChainID. Expected %v, got %v",
			currentBlock.Header.ChainID, self.block.Header.ChainID)
	}
	// 2. hash
	if self.block.Header.PrevBlockHash != currentBlock.HeaderHash {
		return fmt.Errorf("Wrong Block.Header.PrevBlockHash. Expected %v, got %v",
			currentBlock.HeaderHash, self.block.Header.PrevBlockHash)
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
	//5. header hash
	var defaultHash types.Hash
	if !bytes.Equal(defaultHash[:], self.block.HeaderHash[:]) {
		headerHash := vcommon.HeaderHash(self.block)
		if self.block.HeaderHash != headerHash {
			return fmt.Errorf("wrong Block.HeaderHash. Expected %v, got %v",
				headerHash, self.block.HeaderHash)
		}
	}

	var (
		receipts types.Receipts
		allLogs  []*types.Log
		gp       = new(common.GasPool).AddGas(uint64(65536))
	)
	// 6. verify every transactions in the block by evm
	for i, tx := range self.block.Transactions {
		self.chain.Prepare(vcommon.TxHash(tx), vcommon.HeaderHash(self.block), i)
		receipt, _, err := self.VerifyTransaction(self.block.Header.CoinBase, gp, self.block.Header, tx, new(uint64))
		if err != nil {
			log.Error("Tx %x verify failed with error %v.", vcommon.TxHash(tx), err)
			return err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}
	receiptsHash := make([]types.Hash, 0, len(receipts))
	for _, t := range receipts {
		receiptsHash = append(receiptsHash, common.ReceiptHash(t))
		log.Info("Record tx's %x receipt with hash %x.", t.TxHash, common.ReceiptHash(t))
	}
	receiptHash := merkle_tree.ComputeMerkleRoot(receiptsHash)
	var tempHash types.Hash
	if !bytes.Equal(tempHash[:], self.block.Header.ReceiptsRoot[:]) {
		log.Warn("Receipts root has assigned with %x.", self.block.Header.ReceiptsRoot)
		if !bytes.Equal(receiptHash[:], self.block.Header.ReceiptsRoot[:]) {
			log.Error("Receipts root has assigned with %x, but not consistent with %x.",
				self.block.Header.ReceiptsRoot, receiptHash)
			return fmt.Errorf("receipts hash not consistent")
		}
	} else {
		log.Info("Assign receipts hash %x to block %d.", receiptHash, self.block.Header.Height)
		self.block.Header.ReceiptsRoot = receiptHash
	}
	// 7. verify digest if it exists
	if !bytes.Equal(defaultHash[:], self.block.Header.MixDigest[:]) {
		digestHash := vcommon.HeaderDigest(self.block.Header)
		if !bytes.Equal(digestHash[:], self.block.Header.MixDigest[:]) {
			log.Error("Block digest not consistent which assignment is [%x], while compute is [%x].",
				self.block.Header.MixDigest, digestHash)
			return fmt.Errorf("digest not consistent")
		}
	}

	self.receipts = receipts
	self.logs = allLogs

	return nil
}

func (self *Worker) VerifyTransaction(author types.Address, gp *common.GasPool, header *types.Header,
	tx *types.Transaction, usedGas *uint64) (*types.Receipt, uint64, error) {
	if self.VerifyTrsSignature(tx) == false {
		log.Error("Transaction signature verify failed.")
		return nil, 0, fmt.Errorf("transaction signature failed")
	}
	context := evm.NewEVMContext(*tx, header, self.chain, author)
	evmEnv := evm.NewEVM(context, self.chain)
	_, gas, failed, err := ApplyTransaction(evmEnv, tx, gp)
	if err != nil {
		log.Error("Apply transaction %x failed with error %v.", vcommon.TxHash(tx), err)
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
		receipt.ContractAddress = crypto.CreateAddress(evmEnv.Context.Origin, tx.Data.AccountNonce)
		log.Info("Contract address of tx [%x] is [%v]", receipt.TxHash, receipt.ContractAddress)
	}
	// Set the receipt logs and create a bloom for filtering
	// receipt.Logs = self.chain.GetLogs(tx.Hash())
	// receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}

func (self *Worker) VerifyTrsSignature(tx *types.Transaction) bool {
	signer := new(wallett.FrontierSigner)
	from, err := wallett.Sender(signer, tx)
	if nil != err {
		log.Error("Get from by tx's %x signer failed with %v.", vcommon.TxHash(tx), err)
		return false
	}
	if !bytes.Equal((*(tx.Data.From))[:], from.Bytes()) {
		log.Error("Transaction signature verify failed, tx.Data.From is %x, while signed from is %x.", *tx.Data.From, from)
		return false
	}
	return true
}

func (self *Worker) GetReceipts() types.Receipts {
	log.Info("Get Receipts of worker.")
	return self.receipts
}
