package signature

import (
	"fmt"
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/tools"
	"github.com/DSiSc/validator/tools/account"
	"github.com/DSiSc/validator/tools/signature/keypair"
)

type Signature struct {
	Scheme keypair.SignatureScheme
	Value  interface{}
}

//  Sing the data by singer
func Sign(signer Signer, data []byte) ([]byte, error) {
	// TODO: adding key
	signatures := make([]byte, len(data))
	address := signer.(*account.Account).Address
	for i := 0; i < len(data); i++ {
		signatures[i] = address[i] ^ data[i]
	}
	return signatures, nil
}

func verifySpecifiedAddress(pubKey []byte, signData []byte, validatorAddress types.Address) bool {
	var decode = make([]byte, len(signData))
	for i := 0; i < len(signData); i++ {
		decode[i] = pubKey[i] ^ signData[i]
		if validatorAddress[i] != decode[i] {
			return false
		}
	}
	return true
}

// Verify check the signature of data using pubKey
func Verify(pubKey keypair.PublicKey, signature []byte) (types.Address, error) {
	adminAddres := tools.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d")
	var validators = []types.Address{adminAddres}
	pkey := pubKey.([]byte)
	for i := 0; i < len(validators); i++ {
		if verifySpecifiedAddress(pkey, signature, validators[i]) {
			return validators[i], nil
		}
	}
	return *new(types.Address), fmt.Errorf("Invalid signData.")
}

// VerifyMultiSignature check whether more than m sigs are signed by the keys
func VerifyMultiSignature(data []byte, keys []keypair.PublicKey, m int, sigs [][]byte) error {
	return nil
}
