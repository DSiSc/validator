package signature

import (
	"fmt"
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

// Verify check the signature of data using pubKey
func Verify(pubKey keypair.PublicKey, signature []byte) error {
	adminAddres := tools.HexToAddress("333c3310824b7c685133f2bedb2ca4b8b4df633d")
	pkey := pubKey.([]byte)
	var decode = make([]byte, len(signature))
	for i := 0; i < len(signature); i++ {
		decode[i] = pkey[i] ^ signature[i]
		if adminAddres[i] != decode[i] {
			return fmt.Errorf("Signature not consis in pubKey.")
		}
	}
	return nil
}

// VerifyMultiSignature check whether more than m sigs are signed by the keys
func VerifyMultiSignature(data []byte, keys []keypair.PublicKey, m int, sigs [][]byte) error {
	return nil
}
