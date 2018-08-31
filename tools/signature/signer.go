package signature

import (
	"github.com/DSiSc/validator/tools/signature/keypair"
)

// Signer is the abstract interface of user's information(Keys) for signing data.
type Signer interface {
	//get signer's private key
	PrivKey() keypair.PrivateKey

	//get signer's public key
	PubKey() keypair.PublicKey

	Scheme() keypair.SignatureScheme
}
