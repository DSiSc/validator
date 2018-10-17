package worker

import (
	"github.com/DSiSc/blockchain"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	evmNg "github.com/DSiSc/evm-NG"
	vcommon "github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/worker/common"
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
	state      *blockchain.BlockChain
	evm        *evmNg.EVM
	from       types.Address
	to         types.Address
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *evmNg.EVM, trx *types.Transaction, gp *common.GasPool) *StateTransition {
	var receive types.Address
	if trx.Data.Recipient == nil /* contract creation */ {
		receive = types.Address{}
	} else {
		receive = *trx.Data.Recipient
	}
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		tx:       trx,
		from:     *trx.Data.From,
		to:       receive,
		gasPrice: trx.Data.Price,
		value:    trx.Data.Amount,
		data:     trx.Data.Payload,
		state:    evm.StateDB,
		gas:      trx.Data.GasLimit,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyTransaction(evm *evmNg.EVM, tx *types.Transaction, gp *common.GasPool) ([]byte, uint64, bool, error) {
	return NewStateTransition(evm, tx, gp).TransitionDb()
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the used gas. It returns an error if failed.
// An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, usedGas uint64, failed bool, err error) {
	from := *st.tx.Data.From
	sender := vcommon.NewRefAddress(from)
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
	var (
		evm   = st.evm
		vmerr error
	)
	if contractCreation {
		// ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)
		ret, _, st.gas, vmerr = evm.Create(sender, st.data, math.MaxUint64, st.value, st.tx.Data.AccountNonce)
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(from, st.tx.Data.AccountNonce)
		ret, st.gas, vmerr = evm.Call(sender, st.to, st.data, math.MaxUint64, st.value)
	}
	if vmerr != nil {
		log.Error("VM returned with error %v.", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == evmNg.ErrInsufficientBalance {
			return ret, 0, false, vmerr
		}
	}

	return ret, st.gasUsed(), vmerr != nil, vmerr
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
