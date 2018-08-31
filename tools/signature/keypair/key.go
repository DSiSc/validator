package keypair

import (
	"crypto"
)

type SignatureScheme byte

type PublicKey crypto.PublicKey

type PrivateKey interface {
	crypto.PrivateKey
	Public() crypto.PublicKey
}
