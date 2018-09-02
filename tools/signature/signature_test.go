package signature

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Sign(t *testing.T) {
	assert := assert.New(t)
	signature, err := Sign(nil, nil)
	assert.NotNil(signature)
	assert.Nil(err)
	var exceptSign = []byte{
		0x33, 0x3c, 0x33, 0x10, 0x82, 0x4b, 0x7c, 0x68, 0x51, 0x33,
		0xf2, 0xbe, 0xdb, 0x2c, 0xa4, 0xb8, 0xb4, 0xdf, 0x63, 0x3d,
	}
	assert.Equal(exceptSign, signature, "They should be equal")
}

func Test_Verify(t *testing.T) {
	assert := assert.New(t)
	err := Verify(nil, nil, nil)
	assert.Nil(err)
}

func Test_VerifyMultiSignature(t *testing.T) {
	assert := assert.New(t)
	err := VerifyMultiSignature(nil, nil, 0, nil)
	assert.Nil(err)
}
