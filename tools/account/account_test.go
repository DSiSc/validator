package account

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/tools"
	"github.com/DSiSc/validator/tools/signature/keypair"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Account(t *testing.T) {
	assert := assert.New(t)

	var address = types.Address{
		0x33, 0x3c, 0x33, 0x10, 0x82, 0x4b, 0x7c, 0x68, 0x51, 0x33,
		0xf2, 0xbe, 0xdb, 0x2c, 0xa4, 0xb8, 0xb4, 0xdf, 0x63, 0x3d,
	}

	var account = &Account{
		PrivateKey: nil,
		PublicKey:  nil,
		Address:    tools.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d"),
		SigScheme:  keypair.SignatureScheme('a'),
	}

	priKey := account.PrivKey()
	assert.Nil(priKey)

	pubKey := account.PubKey()
	assert.Nil(pubKey)

	scheme := account.Scheme()
	assert.Equal(keypair.SignatureScheme('a'), scheme)

	addr := account.Address
	assert.Equal(address, addr)

	exceptStr := "3<3\x10\x82K|hQ3\xf2\xbe\xdb,\xa4\xb8\xb4\xdfc="
	str := string(address[:])
	assert.Equal(exceptStr, str)

	data := []byte(str)
	assert.Equal(address, tools.BytesToAddress(data))
}
