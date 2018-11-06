package account

import (
	"github.com/DSiSc/craft/types"
	"github.com/DSiSc/validator/tools/signature/keypair"
)

type Account struct {
	PrivateKey keypair.PrivateKey
	PublicKey  keypair.PublicKey
	Address    types.Address
	SigScheme  keypair.SignatureScheme
	Extension  AccountExtension
}

func (*Account) PrivKey() keypair.PrivateKey {
	return nil
}

//get signer's public key
func (*Account) PubKey() keypair.PublicKey {
	return nil
}

func (*Account) Scheme() keypair.SignatureScheme {
	//var temp keypair.SignatureScheme
	var byt byte = 'a'
	return keypair.SignatureScheme(byt)
}

type AccountExtension struct {
	Id  uint64
	Url string
}
