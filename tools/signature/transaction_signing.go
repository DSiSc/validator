package signature

import (
	"github.com/DSiSc/craft/types"
	"math/big"
)

type TxSigner interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *types.Transaction) (types.Address, error)

	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *types.Transaction, sig []byte) (r, s, v *big.Int, err error)

	// Hash returns the hash to be signed.
	Hash(tx *types.Transaction) types.Hash

	// Equal returns true if the given signer is the same as the receiver.
	Equal(TxSigner) bool
}

type sigCache struct {
	signer TxSigner
	from   types.Address
}

func Sender(signer TxSigner, tx *types.Transaction) (types.Address, error) {
	if sc := tx.From.Load(); sc != nil {
		sigCache := sc.(sigCache)
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return types.Address{}, err
	}
	tx.From.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}
