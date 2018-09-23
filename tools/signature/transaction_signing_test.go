package signature

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/common"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type MockSigner struct {
	r     big.Int
	s     big.Int
	v     big.Int
	equal bool
}

func NewMockSigner(r big.Int, s big.Int, v big.Int, equal bool) *MockSigner {
	return &MockSigner{
		r:     r,
		s:     s,
		v:     v,
		equal: equal,
	}
}

var from = types.Address{
	0x12, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
}

func (*MockSigner) Sender(tx *types.Transaction) (types.Address, error) {
	return from, nil
}

func (self *MockSigner) SignatureValues(tx *types.Transaction, sig []byte) (r, s, v *big.Int, err error) {
	r, s, v = &self.r, &self.s, &self.v
	err = nil
	return
}

func (*MockSigner) Hash(tx *types.Transaction) types.Hash {
	return common.TxHash(tx)
}

func (self *MockSigner) Equal(TxSigner) bool {
	return self.equal
}

func TestSender(t *testing.T) {
	to := &types.Address{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	data := []byte{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
	}
	tx := &types.Transaction{
		Data: types.TxData{
			AccountNonce: 0,
			Price:        new(big.Int).SetUint64(10),
			GasLimit:     100,
			Recipient:    to,
			From:         nil,
			Amount:       new(big.Int).SetUint64(50),
			Payload:      data,
		},
	}
	signer := NewMockSigner(*new(big.Int).SetUint64(1), *new(big.Int).SetUint64(2), *new(big.Int).SetUint64(3), true)
	addr, err := Sender(signer, tx)
	assert.Nil(t, err)
	assert.Equal(t, addr, from)

	tx.Data.From = &from
	signer = NewMockSigner(*new(big.Int).SetUint64(1), *new(big.Int).SetUint64(2), *new(big.Int).SetUint64(3), true)
	addr, err = Sender(signer, tx)
	assert.Nil(t, err)
	assert.Equal(t, addr, from)
}
