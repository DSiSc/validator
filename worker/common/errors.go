package common

import (
	"errors"
)

var (
	// ErrGasLimitReached is returned by the gas pool if the amount of gas required
	// by a transaction is higher than what's left in the block.
	ErrGasLimitReached           = errors.New("gas reached limit")
	ErrInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
)
