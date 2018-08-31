package signature

import (
	"github.com/DSiSc/validator/tools/signature/keypair"
)

type Signature struct {
	Scheme keypair.SignatureScheme
	Value  interface{}
}

func Sign(signer Signer, data []byte) ([]byte, error) {
	// TODO: adding key, we use node id to instand
	var sign []byte = []byte("solo_node")
	return sign, nil
}

// Verify check the signature of data using pubKey
func Verify(pubKey keypair.PublicKey, data, signature []byte) error {
	return nil
}

// VerifyMultiSignature check whether more than m sigs are signed by the keys
func VerifyMultiSignature(data []byte, keys []keypair.PublicKey, m int, sigs [][]byte) error {
	return nil
}
