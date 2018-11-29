package signature

import (
	"fmt"
	"github.com/DSiSc/craft/log"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools"
	"github.com/DSiSc/validator/tools/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

var MockAccount = &account.Account{
	Address: types.Address{0x33, 0x3c, 0x33, 0x10, 0x82, 0x4b, 0x7c, 0x68, 0x51, 0x33, 0xf2, 0xbe, 0xdb, 0x2c, 0xa4, 0xb8, 0xb4, 0xdf, 0x63, 0x3d},
}

var MockAccount1 = &account.Account{
	Address: types.Address{0x34, 0x3c, 0x33, 0x10, 0x82, 0x4b, 0x7c, 0x68, 0x51, 0x33, 0xf2, 0xbe, 0xdb, 0x2c, 0xa4, 0xb8, 0xb4, 0xdf, 0x63, 0x3d},
}

var MockSignData = []byte{
	0x8a, 0x73, 0x60, 0x64, 0x7e, 0xae, 0x91, 0xd4, 0xdf, 0x19,
	0x74, 0x29, 0x1a, 0x7f, 0x95, 0xdf, 0xca, 0xb1, 0xdc, 0x36,
}

var Digest = tools.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")

func Test_Sign(t *testing.T) {
	assert := assert.New(t)
	exc, err := Sign(MockAccount, Digest[:])
	assert.Nil(err)
	assert.Equal(MockSignData, exc)
}

func TestVerifyMultiSignature(t *testing.T) {
	assert := assert.New(t)
	signature, err := Sign(MockAccount, Digest[:])
	assert.Nil(err)
	signer, err := Verify(common.ByteToHash(Digest[:]), signature)
	assert.Nil(err)
	assert.Equal(MockAccount.Address, signer)
}

func Test_Verify(t *testing.T) {
	assert := assert.New(t)
	signer, err := Verify(common.ByteToHash(MockSignData), Digest[:])
	assert.Nil(err)
	assert.Equal(signer, tools.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d"))
	NoSignDigest := tools.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf00")
	_, err = Verify(common.ByteToHash(MockSignData), NoSignDigest[:])
	expect := fmt.Errorf("invalid signData")
	assert.Equal(expect, err)

	digest := tools.FromHex("f3b6a8c32257a09ce1b510234fed018212e2910fbad95ba42b1e9e333932bf3f")
	var MockAccount1 = &account.Account{
		Address: types.Address{0x33, 0x3c, 0x33, 0x10, 0x82, 0x4b, 0x7c, 0x68, 0x51, 0x33, 0xf2, 0xbe, 0xdb, 0x2c, 0xa4, 0xb8, 0xb4, 0xdf, 0x63, 0x3d},
	}
	signature, err := Sign(MockAccount1, digest)
	var b types.Hash
	setBytes(digest, &b)
	address, err := Verify(b, signature)
	assert.Nil(err)
	assert.NotNil(address)
	assert.Equal(MockAccount1.Address, address)
}

func Test_VerifyMultiSignature(t *testing.T) {
	assert := assert.New(t)
	err := VerifyMultiSignature(nil, nil, 0, nil)
	assert.Nil(err)
}

func signDataVerify(account account.Account, sign []byte, digest types.Hash) bool {
	address, err := Verify(digest, sign)
	if nil != err {
		log.Error("verify sign %v failed with err %s", sign, err)
	}
	return account.Address == address
}

func setBytes(b []byte, a *types.Hash) {
	copy(a[:], b[:types.HashLength])
}
