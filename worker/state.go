package worker

import (
	"errors"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	evmNg "github.com/DSiSc/evm-NG"
	"github.com/DSiSc/repository"
	evmCommon "github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/worker/common"
	wasmExec "github.com/DSiSc/wasm/exec"
	wasmModule "github.com/DSiSc/wasm/wasm"
	"math"
	"math/big"
)

type StateTransition struct {
	gp         *common.GasPool
	tx         *types.Transaction
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      *repository.Repository
	from       types.Address
	to         types.Address
	nonce      uint64
	header     *types.Header
	author     types.Address
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(author types.Address, header *types.Header, chain *repository.Repository, trx *types.Transaction, gp *common.GasPool) *StateTransition {
	var receive types.Address
	if trx.Data.Recipient == nil /* contract creation */ {
		receive = types.Address{}
	} else {
		receive = *trx.Data.Recipient
	}
	return &StateTransition{
		author:     author,
		gp:         gp,
		tx:         trx,
		from:       *trx.Data.From,
		to:         receive,
		gasPrice:   trx.Data.Price,
		value:      trx.Data.Amount,
		data:       trx.Data.Payload,
		state:      chain,
		gas:        trx.Data.GasLimit,
		initialGas: math.MaxUint64,
		nonce:      trx.Data.AccountNonce,
		header:     header,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyTransaction(author types.Address, header *types.Header, chain *repository.Repository, tx *types.Transaction, gp *common.GasPool) ([]byte, uint64, bool, error, types.Address) {
	return NewStateTransition(author, header, chain, tx, gp).TransitionDb()
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the used gas. It returns an error if failed.
// An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, usedGas uint64, failed bool, err error, address types.Address) {
	if err = st.preCheck(); err != nil {
		return
	}
	if st.isWasmContract(st.tx) {
		ret, address, st.gas, err = st.execWasmContract()
	} else {
		var vmerr error
		ret, address, st.gas, vmerr = st.execSolidityContract()
		if vmerr != nil {
			log.Debug("VM returned with error", "err", vmerr)
			// The only possible consensus-error would be if there wasn't
			// sufficient balance to make the transfer happen. The first
			// balance transfer may never fail.
			if vmerr == evmNg.ErrInsufficientBalance {
				return ret, 0, false, vmerr, address
			}
		}
	}
	return ret, st.gasUsed(), err != nil, err, address
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(*st.tx.Data.From, remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

// check tx's nonce
func (st *StateTransition) preCheck() error {
	// Make sure this transaction's nonce is correct.
	nonce := st.state.GetNonce(st.from)
	if nonce < st.nonce {
		return errors.New("blacklisted hash")
	} else if nonce > st.nonce {
		return errors.New("nonce too high")
	}
	return nil
}

// check whether the tx Recipient is wasm contract
func (st *StateTransition) isWasmContract(tx *types.Transaction) bool {
	var code []byte
	if (nil == tx.Data.Recipient || types.Address{} == *tx.Data.Recipient) {
		code = tx.Data.Payload
	} else {
		code = st.state.GetCode(*tx.Data.Recipient)
	}
	return wasmModule.IsValidWasmCode(code)
}

func (st *StateTransition) execSolidityContract() (ret []byte, contractAddr types.Address, leftOverGas uint64, err error) {
	from := *st.tx.Data.From
	sender := evmCommon.NewRefAddress(from)
	//homestead := st.evm.ChainConfig().IsHomestead(st.evm.BlockNumber)
	contractCreation := st.tx.Data.Recipient == nil
	/*
		// Pay intrinsic gas
		gas, err := IntrinsicGas(st.data, contractCreation, homestead)
		if err != nil {
			return nil, 0, false, err
		}
		if err = st.useGas(gas); err != nil {
			return nil, 0, false, err
		}
	*/
	context := evmNg.NewEVMContext(*st.tx, st.header, st.state, st.author)
	evm := evmNg.NewEVM(context, st.state)

	if contractCreation {
		// ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)
		ret, contractAddr, leftOverGas, err = evm.Create(sender, st.data, math.MaxUint64, st.value)
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(from, st.state.GetNonce(sender.Address())+1)
		ret, leftOverGas, err = evm.Call(sender, st.to, st.data, math.MaxUint64, st.value)
	}
	return ret, contractAddr, leftOverGas, err
}

func (st *StateTransition) execWasmContract() (ret []byte, contractAddr types.Address, leftOverGas uint64, err error) {
	context := wasmExec.NewWasmChainContext(st.tx, st.header, st.state, st.author)
	wvm := wasmExec.NewVM(context, st.state)
	contractCreation := st.tx.Data.Recipient == nil
	if contractCreation {
		ret, contractAddr, leftOverGas, err = wvm.Create(*st.tx.Data.From, st.data, math.MaxUint64, st.value)
	} else {
		st.state.SetNonce(*st.tx.Data.From, st.state.GetNonce(*st.tx.Data.From)+1)
		ret, leftOverGas, err = wvm.Call(*st.tx.Data.From, *st.tx.Data.Recipient, st.data, math.MaxUint64, st.value)
	}
	return ret, contractAddr, leftOverGas, err
}
