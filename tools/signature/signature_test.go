package signature

import (
	"fmt"
	"github.com/DSiSc/validator/common"
	"github.com/DSiSc/validator/tools"
	"github.com/DSiSc/validator/tools/account"
	"github.com/stretchr/testify/assert"
	"testing"
)

var MockAccount = &account.Account{
	Address: tools.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d"),
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

func Test_Verify(t *testing.T) {
	assert := assert.New(t)
	signer, err := Verify(common.ByteToHash(MockSignData), Digest[:])
	assert.Nil(err)
	assert.Equal(signer, tools.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d"))
	NoSignDigest := tools.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf00")
	_, err = Verify(common.ByteToHash(MockSignData), NoSignDigest[:])
	expect := fmt.Errorf("Invalid signData.")
	assert.Equal(expect, err)
}

func Test_VerifyMultiSignature(t *testing.T) {
	assert := assert.New(t)
	err := VerifyMultiSignature(nil, nil, 0, nil)
	assert.Nil(err)
}
