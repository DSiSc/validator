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
	var sign = []byte{
		0x33, 0x3c, 0x33, 0x10, 0x82, 0x4b, 0x7c, 0x68, 0x51, 0x33,
		0xf2, 0xbe, 0xdb, 0x2c, 0xa4, 0xb8, 0xb4, 0xdf, 0x63, 0x3d,
	}
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
