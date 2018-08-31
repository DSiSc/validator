package account

import (
	"github.com/DSiSc/validator/tools"
	"github.com/DSiSc/validator/tools/signature/keypair"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Account(t *testing.T) {
	assert := assert.New(t)

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

}
